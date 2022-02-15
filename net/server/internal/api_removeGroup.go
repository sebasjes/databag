package databag

import (
  "errors"
  "net/http"
  "gorm.io/gorm"
  "github.com/gorilla/mux"
  "databag/internal/store"
)

func RemoveGroup(w http.ResponseWriter, r *http.Request) {

  account, code, err := BearerAppToken(r, false);
  if err != nil {
    ErrResponse(w, code, err)
    return
  }

  // scan parameters
  params := mux.Vars(r)
  groupId := params["groupId"]

  // load referenced group
  var slot store.GroupSlot
  if err := store.DB.Preload("Group.Cards.CardSlot").Where("account_id = ? AND group_slot_id = ?", account.ID, groupId).First(&slot).Error; err != nil {
    if !errors.Is(err, gorm.ErrRecordNotFound) {
      ErrResponse(w, http.StatusInternalServerError, err)
    } else {
      ErrResponse(w, http.StatusNotFound, err)
    }
    return
  }
  if slot.Group == nil {
    ErrResponse(w, http.StatusNotFound, errors.New("group has been deleted"))
    return
  }

  var cards []*store.Card
  for _, card := range slot.Group.Cards {
    cards = append(cards, &card)
  }

  // save and update contact revision
  err = store.DB.Transaction(func(tx *gorm.DB) error {
    if res := tx.Model(&slot.Group).Association("Channels").Clear(); res != nil {
      return res
    }
    if res := tx.Model(&slot.Group).Association("Articles").Clear(); res != nil {
      return res
    }
    for _, card := range cards {
      if res := tx.Model(card).Update("view_revision", card.ViewRevision + 1).Error; res != nil {
        return res
      }
      if res := tx.Model(card).Update("detail_revision", card.DetailRevision + 1).Error; res != nil {
        return res
      }
      if res := tx.Model(card.CardSlot).Update("revision", account.CardRevision + 1).Error; res != nil {
        return res
      }
    }
    if res := tx.Model(&slot.Group).Association("Cards").Clear(); res != nil {
      return res
    }
    if res := tx.Delete(&slot.Group).Error; res != nil {
      return res
    }
    slot.Group = nil
    if res := tx.Model(&slot).Update("group_id", 0).Error; res != nil {
      return res
    }
    if res := tx.Model(&slot).Update("revision", account.GroupRevision + 1).Error; res != nil {
      return res
    }
    if res := tx.Model(&account).Update("group_revision", account.GroupRevision + 1).Error; res != nil {
      return res
    }
    if res := tx.Model(&account).Update("card_revision", account.CardRevision + 1).Error; res != nil {
      return res
    }
    return nil
  })
  if err != nil {
    ErrResponse(w, http.StatusInternalServerError, err)
    return
  }

  for _, card := range cards {
    SetContactViewNotification(account, card)
  }
  SetStatus(account)
  WriteResponse(w, nil);
}

package databag

import (
  "time"
  "bytes"
  "errors"
  "net/http"
  "gorm.io/gorm"
  "databag/internal/store"
  "encoding/base64"
)

func GetAccountListingImage(w http.ResponseWriter, r *http.Request) {

  var account store.Account
  if err := store.DB.Preload("AccountDetail").Where("searchable = ? AND disabled = ?", true, false).First(&account).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      ErrResponse(w, http.StatusNotFound, err)
    } else {
      ErrResponse(w, http.StatusInternalServerError, err)
    }
    return
  }

  if account.AccountDetail.Image == "" {
    ErrResponse(w, http.StatusNotFound, errors.New("image not set"))
    return
  }

  data, err := base64.StdEncoding.DecodeString(account.AccountDetail.Image)
  if err != nil {
    ErrResponse(w, http.StatusNotFound, errors.New("image not valid"))
    return
  }

  // response with content
  http.ServeContent(w, r, "image", time.Unix(account.Updated, 0), bytes.NewReader(data))
}


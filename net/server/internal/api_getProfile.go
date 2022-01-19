package databag

import (
  "net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {

  account, res := BearerAppToken(r, true);
  if res != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  if account.Disabled {
    w.WriteHeader(http.StatusGone);
    return
  }
  detail := account.AccountDetail

  // send profile data
  profile := Profile {
    Guid: account.Guid,
    Handle: account.Username,
    Name: detail.Name,
    Description: detail.Description,
    Location: detail.Location,
    Image: detail.Image,
    Revision: account.ProfileRevision,
    Version: APP_VERSION,
    Node: "https://" + getStrConfigValue(CONFIG_DOMAIN, "") + "/",
  }
  WriteResponse(w, profile)
}

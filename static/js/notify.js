function notify(msg, msgType) {
    notie.alert({
    type: msgType, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg,
})
  }

  function alertError(msg) {
    notie.alert({
    type: "error", // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg,
})
  }

  function alertWarning(msg) {
    notie.alert({
    type: "warning", // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg,
})
  }
  
  {{with .Error}}
  notify("{{.}}", "error")
  {{end}}

  {{with .Flash}}
  notify("{{.}}", "success")
  {{end}}

  {{with .Warning}}
  notify("{{.}}", "warning")
  {{end}}

 
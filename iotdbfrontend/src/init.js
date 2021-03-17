fetch(`/csrf`, {
    method: 'GET',
    credentials: "same-origin"
})
    .catch(alert("Error: Server Unavailable. Please reload"))


const initCSRF = async () => {
    try {
        const response = await fetch(`/csrf`, {
            method: 'GET',
            credentials: "same-origin"
        })

        if (response.status != 200) {
            response.text().then(text => alert(text))
        }
    } catch (e) {
        alert("Error: Server Unavailable. Please try again")
    }
}

export default initCSRF;

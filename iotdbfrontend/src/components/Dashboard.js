import React, { useState } from 'react';
import LogInForm from './LogInForm';
import LogOutButton from './LogOutButton'
import Cookies from 'js-cookie';


const Dashboard = () => {
    var jwtExists = false
    console.log(Cookies.get('logged_in'))
    if (Cookies.get('logged_in') == 'true') {
        jwtExists = true
    }
    const [loggedIn, setLoggedIn] = useState(jwtExists)

    const handleLogout = async (event) => {
        // Added preventDetfault to stop the page from reloading before the response returns
        event.preventDefault();
        try {
            const requestOptions = {
                method: 'POST',
                body: JSON.stringify({ csrf: Cookies.get('CSRF') })
            };

            const response = await fetch(`/logout`, requestOptions)

            if (response.status != 200) {
                response.text().then(text => alert(text))
            }
        }
        catch (e) {
            alert("Error: Server Unavailable. Please try again")
        }
        //this is placed outside of the try/catch so that errors force the user to reauthenticate
        Cookies.remove('JWT')
        Cookies.remove('logged_in')
        setLoggedIn(false)
    }

    if (loggedIn) {
        return (
            <header className="top-nav">
                <h1>
                    User Management Dashboard
                </h1>
                <LogOutButton handleClick={handleLogout} />
            </header>
        )
    }
    return <LogInForm setLoggedIn={setLoggedIn} />
}


export default Dashboard;

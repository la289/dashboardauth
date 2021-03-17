import React, { useState } from 'react';
import LogInForm from './LogInForm';
import LogOutButton from './LogOutButton';
import Cookies from 'js-cookie';
import { authState } from '../constants.js';

const Dashboard = () => {
    const [loggedIn, setLoggedIn] = useState(Cookies.get('logged_in') == 'true')

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
        Cookies.remove(authState.LOGGED_IN)
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

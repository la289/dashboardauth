import React, {useState} from 'react';
import LogInForm from './LogInForm';
import LogOutButton from './LogOutButton'
import Cookies from 'js-cookie';


const Dashboard = () => {
    const [loggedIn, setLoggedIn] = useState(false)

    const handleLogout = async (event) => {
        // Added preventDetfault to stop the page from reloading before the response returns
        event.preventDefault();
        try{
            const requestOptions = {
                method: 'POST',
                body: JSON.stringify({csrf: Cookies.get('CSRF')})
            };

            const response = await fetch(`/logout`, requestOptions);

            if (response.status == 200){
                setLoggedIn(false);
            } else {
                response.text().then(text => alert(text));
            }
        }
        catch (e) {
            console.log(e)
            alert("Server Unavailable")
        }
    }

    if (loggedIn) {
        return (
            <header class="top-nav">
            <h1>
              User Management Dashboard
            </h1>
            <button class="button is-border" onClick = {handleLogout}>Logout</button>
          </header>
        )
    }

    return <LogInForm setLoggedIn = {setLoggedIn}/>
}


export default Dashboard;

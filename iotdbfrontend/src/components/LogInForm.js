import React, { useState } from 'react';
import LogInButton from './LogInButton.js';
import FormTextField from './FormTextField.js';
import Cookies from 'js-cookie';
import { authState } from '../constants.js';


const LogInForm = ({ setLoggedIn }) => {

    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    const handleClick = async (event) => {
        // Added preventDetfault to stop the page from reloading before the response returns
        event.preventDefault();
        try {
            const requestOptions = {
                method: 'POST',
                body: JSON.stringify({ email: email, password: password, csrf: Cookies.get('CSRF') })
            };

            const response = await fetch(`/login`, requestOptions);

            if (response.status == 200) {
                setLoggedIn(true);
                //need new cookie since JWT is httponly
                Cookies.set(authState.LOGGED_IN, 'true')

            } else {
                response.text().then(text => alert(text));
            }
        }
        catch (e) {
            alert("Error: Server Unavailable. Please try again")
        }
    }

    return (
        <form className="login-form" onSubmit={handleClick}>
            <h1>Sign Into Your Account</h1>

            <FormTextField name="email" onChange={setEmail} />

            <FormTextField name="password" onChange={setPassword} />

            <LogInButton />
        </form>
    )


}



export default LogInForm;

import React from 'react';

const LogOutButton = ({ handleClick }) => {
    return (
        <button type="button" className="button is-border" onClick={(e) => handleClick(e)}>Logout</button>
    )
}

export default LogOutButton;

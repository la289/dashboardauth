import React from 'react';

 const LogInButton = ({handleClick}) => {
    return (
    <input type="submit" value="Login to my Dashboard" className="button block" onClick={handleClick} />
    )
}

export default LogInButton

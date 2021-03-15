import React from 'react';

const FormTextField = ({ name, onChange }) => {
    return (
        <div>
            <label htmlFor={name}>{name}</label>
            <input type={name} id={name} className="field" onChange={e => onChange(e.target.value)} />
        </div>
    )
}

export default FormTextField;

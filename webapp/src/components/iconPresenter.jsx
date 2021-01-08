import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

const Icon = require('/public/Images/icon.svg');

const IconPresenter = ({encryptionEnabled}) => {
    const style = {
        position: 'relative',
        top: '5px',
    };
    return (
        <FormattedMessage
            id='anonymous.encryption.ariaLabel'
            defaultMessage='encryption icon'
        >
            {(ariaLabel) => (
                <span
                    style={style}
                    aria-label={ariaLabel}
                >
                    <img
                        src={Icon}
                        style={{fill: encryptionEnabled ? '' : '#2389d7'}}
                        alt='Image not loaded'
                    />
                </span>
            )}
        </FormattedMessage>
    );
};

IconPresenter.propTypes = {
    encryptionEnabled: PropTypes.bool.isRequired,
};

export default IconPresenter;

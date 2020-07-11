import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';

import SVGs from '../constants/SVGs';

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
                    dangerouslySetInnerHTML={{__html: encryptionEnabled ? SVGs.ANONYMOUS_ICON_ENABLED : SVGs.ANONYMOUS_ICON_DISABLED}}
                />
            )}
        </FormattedMessage>
    );
};

IconPresenter.propTypes = {
    encryptionEnabled: PropTypes.bool.isRequired,
};

export default IconPresenter;

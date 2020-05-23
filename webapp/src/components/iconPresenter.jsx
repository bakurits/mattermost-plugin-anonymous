import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import SVGs from '../constants/SVGs';

const IconPresenter = ({encryptionEnabled}) => {
    const style = getStyle();
    return (
        <FormattedMessage
            id='anonymous.encryption.ariaLabel'
            defaultMessage='encryption icon'
        >
            {(ariaLabel) => (
                <span
                    style={style.iconStyle}
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

const getStyle = () => {
    return makeStyleFromTheme(() => {
        return {
            iconStyle: {
                position: 'relative',
                top: '2px',
            },
        };
    });
};

export default IconPresenter;

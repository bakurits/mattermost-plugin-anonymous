import {connect} from 'react-redux';

import Constants from '../constants';

import IconPresenter from './iconPresenter';

const mapStateToProps = (state) => {
    return {
        encryptionEnabled: state[Constants.REDUCER_ID].encryptionState,
    };
};

const Icon = connect(mapStateToProps)(IconPresenter);

export default Icon;

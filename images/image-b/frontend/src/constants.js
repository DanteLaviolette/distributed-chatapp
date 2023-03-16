import PropTypes from 'prop-types';

const constants = {
    PASSWORD_REQUIREMENTS: "Password must be at least 8 characters.",
    AUTH_HEADER: "authorization",
    USER_PROP_TYPE: PropTypes.shape({
        data: PropTypes.shape({
            id: PropTypes.string,
            name: PropTypes.string,
            email: PropTypes.string,
        })
    })
}
export default constants
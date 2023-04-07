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
    }),
    TOP_BAR_HEIGHT: 50,
    MESSAGE_BAR_HEIGHT: "55px",
    TOAST_CONFIG: {
        position: "top-left",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined,
        theme: "colored",
    },
    MS_TO_NS: 1000000
}
export default constants
import requester from "@/utils/request";
import store from "@/store/index";
import Cookies from "js-cookie";
import ChangePassword from "./change-password.vue";

const authService = {
  login: (username, password) => {
    return requester
      .post("/auth/login", {
        email: username,
        password: password,
      })
      .then((response) => {
        const user = response.data;
        store.dispatch("currentUser/setUser", {
          id: user.id_usuario,
          name: user.nome_usuario,
          token: "newToken",
        });
        // Cookies.set("jwtToken", user.token, {
        //   expires: 1, // Optional: set cookie expiration time in days
        //   secure: true, // Optional: ensures the cookie is only sent over HTTPS
        //   sameSite: "Strict", // Optional: prevent CSRF attacks by restricting cross-site access
        // });
        return true;
      })
      .catch((err) => {
        if (err.status === 401) {
          return false;
        }
      });
  },
  logout: () => {
    return store.dispatch("currentUser/setUser", null);
  },
  signup: (user) =>
    requester.post("/auth/signup", user).then((resp) => resp.data),
  confirmAccount: (token) =>
    requester.get(`/auth/verify/${token}`).then((resp) => resp.data),
  startRegistration: (email) =>
    requester
      .post("/auth/register", { email: email })
      .then((resp) => resp.data),
  changePassword: (userData) =>
    requester.post("/auth/changePassword", userData).then((resp) => resp.data),
  redefinePassword: (payload) =>
    requester.post("/auth/redefinePassword", payload).then((resp) => resp.data),
  startRecoverProcess: (email) =>
    requester
      .post("/auth/recoverPassword", { email: email })
      .then((resp) => resp.data),
};

export default authService;

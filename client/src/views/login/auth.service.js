import requester from "@/utils/request";
import store from "@/store/index";
import Cookies from "js-cookie";

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
        });
        Cookies.set("jwtToken", user.token, {
          expires: 1, // Optional: set cookie expiration time in days
          secure: true, // Optional: ensures the cookie is only sent over HTTPS
          sameSite: "Strict", // Optional: prevent CSRF attacks by restricting cross-site access
        });
        return true;
      })
      .catch((err) => {
        if (err.status === 401) {
          return false;
        }
      });
  },
  logout: () => {
    Cookies.remove("jwtToken");
    return store.dispatch("currentUser/setUser", null);
  },
  signup: (user) =>
    requester
      .post("/auth/signup", user)
      .then((resp) => resp.data),
  confirmAccount: (token) =>
    requester.get(`/auth/verify/${token}`).then((resp) => resp.data),
  startRegistration: (email) =>
    requester
      .post(`/auth/register`, { email: email })
      .then((resp) => resp.data),
};

export default authService;

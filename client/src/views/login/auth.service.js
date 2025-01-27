import requester from "@/utils/request";
import store from "@/store/index";

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
          role: user.papel,
          token: user.token,
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
    return store.dispatch("currentUser/setUser", null);
  },
  signup: (user) => {
    return requester.post("/auth/signup", user).then((resp) => {
      return resp.data;
    });
  },
};

export default authService;

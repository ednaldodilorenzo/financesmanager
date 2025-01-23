import axios from "axios";
import { HTTP_STATUS_CODE, SERVER_VALIDATION_CODE } from "../utils/constants";
import router from "../router";
import store from "../store/index";

const request = axios.create({
  baseURL: "/api",
});

request.interceptors.request.use((config) => {
  const token = store.getters["currentUser/getUserToken"];

  return {
    ...config,
    headers: {
      ...config.headers,
      Authorization: `Bearer ${token}`,
    },
  };
});

request.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response.status === HTTP_STATUS_CODE.UNAUTHORIZED) {
      if (
        SERVER_VALIDATION_CODE.INVALID_CREDENTIALS === error.response.data.code
      ) {
        throw error;
      } else {
        localStorage.removeItem("token");
        router.push("/login");
      }
    } else if (
      error.response.status == HTTP_STATUS_CODE.INTERNAL_SERVER_ERROR
    ) {
      router.push("/error");
    } else {
      throw error;
    }
  }
);

const requester = {
  get: (url, params = {}) => {
    const url_params = Object.entries(params)
      .map(([name, value]) => name + "=" + value)
      .join("&");
    return request.get(url_params ? url + "?" + url_params : url);
  },
  post: (url, data) => {
    return request.post(url, data);
  },
  patch: (url, data) => {
    return request.patch(url, data);
  },
  put: (url, data) => {
    return request.put(url, data);
  },
  delete: (url) => {
    return request.delete(url);
  },
};

export default requester;

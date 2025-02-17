import requester from "./request";
import { throttle } from "./support";

class GenericService {
  constructor(serviceURL) {
    this.url = serviceURL;
  }

  findAll = throttle(
    (params = {}) =>
      requester.get(`${this.url}/`, params).then((resp) => resp.data),
    1000
  );

  findById = throttle(
    (id) => requester.get(`${this.url}/${id}`).then((resp) => resp.data),
    1000
  );

  create = throttle(
    (item) => requester.post(`${this.url}/`, item).then((resp) => resp),
    1000
  );

  modify = throttle(
    (id, item) =>
      requester.patch(`${this.url}/${id}`, item).then((resp) => resp.data),
    1000
  );

  delete = throttle(
    (id) => requester.delete(`${this.url}/${id}`).then((resp) => resp.data),
    1000
  );
}

export default GenericService;

import requester from "./request";

class GenericService {
  constructor(serviceURL) {
    this.url = serviceURL;
  }

  findAll(params = {}) {
    return requester.get(`${this.url}/`, params).then((resp) => {
      return resp.data;
    });
  }

  findById(id) {
    return requester.get(`${this.url}/${id}`).then((resp) => {
      return resp.data;
    });
  }

  create(item) {
    return requester.post(`${this.url}/`, item).then((resp) => {
      return resp;
    });
  }

  modify(id, item) {
    return requester.patch(`${this.url}/${id}`, item).then((resp) => {
      return resp.data;
    });
  }

  delete(id) {
    return requester.delete(`${this.url}/${id}`).then((resp) => {
      return resp.data;
    });
  }
}

export default GenericService;

import GenericService from "@/utils/service";
import requester from "@/utils/request";

class TransactionService extends GenericService {
  constructor() {
    // Call the parent constructor with the specific service URL
    super("/transactions");
  }

  // Add a custom method specific to UserService
  prepareForImport(formData) {
    return requester
      .post(`${this.url}/upload`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
      .then((resp) => resp.data);
  }

  sendBatchImport(items) {
    return requester.post(`${this.url}/list`, items).then((resp) => resp.data);
  }

  findWithSplits(transactionId) {
    return requester.get(`${this.url}/${transactionId}/split`).then((resp) => resp.data);
  }

  createSplit(item) {
    return requester.post(`${this.url}/split`, item).then((resp) => resp.data);
  }

  modifySplit(transactionId, item) {
    return requester
      .put(`${this.url}/${transactionId}/split`, item)
      .then((resp) => resp.data);
  }

  getAllByYear(year) {
    return requester.get(`${this.url}/summary/${year}`).then((resp) => resp.data);
  }
}

const transactionService = new TransactionService();

export default transactionService;

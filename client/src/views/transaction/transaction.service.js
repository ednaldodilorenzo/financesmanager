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
}

const transactionService = new TransactionService();

export default transactionService;

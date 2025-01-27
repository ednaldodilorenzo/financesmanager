import { formatCurrency, parseCurrencyToNumber } from "@/utils/numbers";

export default class BudgetItem {
  // Define private field
  strValueChanged = false;

  constructor(item, allCategories) {
    this.id = item.id;
    this.year = item.year;
    this.category = this.findCategory(item.categoryId, allCategories);
    this._value = item.value; // Assume `item.value` is in cents
    this._strValue = formatCurrency("" + item.value);
  }

  // Setter for strValue
  set strValue(value) {
    this._strValue = value;
    this.strValueChanged = true;
  }

  // Getter for strValue
  get strValue() {
    return this._strValue;
  }

  // Getter for numeric value
  get value() {
    if (this.strValueChanged) {
      this._value = parseCurrencyToNumber(this._strValue);
      this.strValueChanged = false;
    }
    return this._value;
  }

  // Method to find the category
  findCategory(categoryId, allCategories) {
    const resultCategory = allCategories.find(
      (category) => category.id === categoryId
    );
    return resultCategory ? resultCategory : { type: "" };
  }

  clone() {
    return Object.create(
      Object.getPrototypeOf(this),
      Object.getOwnPropertyDescriptors(this)
    );
  }

  toJSON() {
    return {
      year: this.year,
      value: this.value,
      categoryId: this.category.id,
      id: this.id,
    };
  }
}

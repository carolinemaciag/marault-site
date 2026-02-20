document.addEventListener("DOMContentLoaded", function () {
  const phoneInput = document.getElementById("phone");
  if (!phoneInput) return;

  phoneInput.addEventListener("input", function (e) {
    let numbers = e.target.value.replace(/\D/g, "");

    if (numbers.length > 10) {
      numbers = numbers.slice(0, 10);
    }

    let formatted = numbers;

    if (numbers.length > 6) {
      formatted = `(${numbers.slice(0,3)}) ${numbers.slice(3,6)}-${numbers.slice(6)}`;
    } else if (numbers.length > 3) {
      formatted = `(${numbers.slice(0,3)}) ${numbers.slice(3)}`;
    } else if (numbers.length > 0) {
      formatted = `(${numbers}`;
    }

    e.target.value = formatted;
  });
});
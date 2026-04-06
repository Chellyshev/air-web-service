document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector(".auth-form");
  const loginInput = document.getElementById("login");
  const passwordInput = document.getElementById("password");

  form.addEventListener("submit", (e) => {
    clearErrors();

    let isValid = true;

    const login = loginInput.value.trim();
    const password = passwordInput.value.trim();

    if (login.length < 3) {
      showError(loginInput, "Минимум 3 символа");
      isValid = false;
    }

    if (password.length < 6) {
      showError(passwordInput, "Минимум 6 символов");
      isValid = false;
    }

    if (!isValid) {
      e.preventDefault(); // блокируем только если ошибка
    }
  });

  function showError(input, message) {
    const error = document.createElement("div");
    error.className = "input-error";
    error.innerText = message;

    input.classList.add("error");
    input.parentNode.appendChild(error);
  }

  function clearErrors() {
    document.querySelectorAll(".input-error").forEach(el => el.remove());
    document.querySelectorAll(".error").forEach(el => el.classList.remove("error"));
  }
});
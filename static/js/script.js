const controlRemoveModal = () => {
  document.addEventListener("htmx:beforeRequest", () => {
    const modal = document.querySelector(".modal");
    modal.nextElementSibling.remove();
    modal.remove();
  });
};

const init = () => {
  controlRemoveModal();
};

document.addEventListener("DOMContentLoaded", () => {
  init();
});

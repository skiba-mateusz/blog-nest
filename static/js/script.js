const controlRemoveModal = () => {
  document.addEventListener("htmx:beforeRequest", () => {
    const modal = document.querySelector(".modal");
    if (!modal) return;
    modal.nextElementSibling.remove();
    modal.remove();
  });
};

const clickOutside = (element, callback) => {
  document.addEventListener("click", (e) => {
    if (!element.contains(e.target)) {
      callback();
    }
  });
};

const controlMenus = () => {
  const menus = document.querySelectorAll(".menu");

  const closeMenu = (menuTrigger, menu) => {
    menuTrigger.setAttribute("aria-expanded", "false");
    menu.classList.remove("menu--active");
  };

  const toggleMenu = (menuTrigger, menu) => {
    const isOpen = menu.classList.toggle("menu--active");
    menuTrigger.setAttribute("aria-expanded", isOpen);

    if (!isOpen) closeMenu(menuTrigger, menu);
  };

  menus.forEach((menu) => {
    const menuTrigger = menu.children[0];

    clickOutside(menu, () => closeMenu(menuTrigger, menu));

    menuTrigger.addEventListener("click", (e) => {
      if (!e.target.closest(".menu__trigger")) return;
      toggleMenu(menuTrigger, menu);
    });
  });
};

const controlHighlightActiveNavLink = () => {
  const navLinks = document.querySelectorAll(".nav__link");
  const currentURL = window.location.href;

  navLinks.forEach((navLink) => {
    currentURL == navLink.href
      ? navLink.classList.add("nav__link--active")
      : navLink.classList.remove("nav__link--active");
  });
};

const init = () => {
  controlHighlightActiveNavLink();
  controlRemoveModal();
  controlMenus();
};

document.addEventListener("DOMContentLoaded", () => {
  init();
});

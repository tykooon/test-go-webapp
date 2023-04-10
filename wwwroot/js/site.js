document.addEventListener("DOMContentLoaded", setActiveNavItem);

function setActiveNavItem() {
    let title = document.querySelector("title").textContent.toLowerCase();
    let itemId = title + "-nav-item";
    let navLink = document.getElementById(itemId);
    navLink.className += " active";
    navLink.setAttribute("aria-current", "page");
}
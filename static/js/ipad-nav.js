/**
 * ipad-nav.js
 * Mirrors mobile-nav.js logic for the iPad layout.
 * Targets .ipad-* classes and body.ipad-base.
 */

(function () {
  "use strict";

  const toggle = document.querySelector(".ipad-nav-toggle");
  const menu = document.querySelector(".ipad-nav-menu");
  const screen = document.querySelector(".ipad-screen");
  const body = document.body;

  if (!toggle || !menu) return;

  // ── Open / close main menu ──────────────────────────────────────────
  function openMenu() {
    toggle.classList.add("is-open");
    menu.classList.add("is-open");
    body.classList.add("ipad-menu-open");
    toggle.setAttribute("aria-expanded", "true");
  }

  function closeMenu() {
    toggle.classList.remove("is-open");
    menu.classList.remove("is-open");
    body.classList.remove("ipad-menu-open");
    toggle.setAttribute("aria-expanded", "false");

    // Also close any open dropdowns
    document.querySelectorAll(".ipad-nav-item.open").forEach(function (item) {
      item.classList.remove("open");
      const caret = item.querySelector(".ipad-nav-caret");
      if (caret) caret.setAttribute("aria-expanded", "false");
    });
  }

  toggle.addEventListener("click", function (e) {
    e.stopPropagation();
    if (menu.classList.contains("is-open")) {
      closeMenu();
    } else {
      openMenu();
    }
  });

  // ── Services dropdown ───────────────────────────────────────────────
  document.querySelectorAll(".ipad-nav-item.has-ipad-dropdown").forEach(function (item) {
    const caret = item.querySelector(".ipad-nav-caret");
    if (!caret) return;

    caret.addEventListener("click", function (e) {
      e.stopPropagation();
      const isOpen = item.classList.contains("open");
      // Close all other dropdowns first
      document.querySelectorAll(".ipad-nav-item.open").forEach(function (other) {
        if (other !== item) {
          other.classList.remove("open");
          const otherCaret = other.querySelector(".ipad-nav-caret");
          if (otherCaret) otherCaret.setAttribute("aria-expanded", "false");
        }
      });
      item.classList.toggle("open", !isOpen);
      caret.setAttribute("aria-expanded", (!isOpen).toString());
    });
  });

  // ── Close menu on outside tap ───────────────────────────────────────
  document.addEventListener("click", function (e) {
    if (!menu.contains(e.target) && !toggle.contains(e.target)) {
      if (menu.classList.contains("is-open")) {
        closeMenu();
      }
    }
  });

  // ── Close menu when a nav link is tapped ───────────────────────────
  menu.querySelectorAll("a:not(.ipad-nav-link)").forEach(function (link) {
    link.addEventListener("click", function () {
      closeMenu();
    });
  });

  // ── Close on Escape ─────────────────────────────────────────────────
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape" && menu.classList.contains("is-open")) {
      closeMenu();
      toggle.focus();
    }
  });

})();
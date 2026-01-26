// =========================
// NAVIGATION DROPDOWN
// =========================

document.addEventListener("DOMContentLoaded", () => {
  const dropdownItems = document.querySelectorAll(".has-dropdown");

  dropdownItems.forEach(item => {
    const caret = item.querySelector(".nav-caret");
    const dropdown = item.querySelector(".dropdown");

    if (!caret || !dropdown) return;

    caret.addEventListener("click", (e) => {
      e.stopPropagation();

      const isOpen = dropdown.classList.contains("open");

      // Close all dropdowns
      document.querySelectorAll(".dropdown.open").forEach(d => {
        d.classList.remove("open");
      });
      document.querySelectorAll(".nav-caret[aria-expanded='true']").forEach(c => {
        c.setAttribute("aria-expanded", "false");
      });

      // Toggle current
      if (!isOpen) {
        dropdown.classList.add("open");
        caret.setAttribute("aria-expanded", "true");
      }
    });
  });

  // Close dropdown when clicking outside
  document.addEventListener("click", () => {
    document.querySelectorAll(".dropdown.open").forEach(d => {
      d.classList.remove("open");
    });
    document.querySelectorAll(".nav-caret[aria-expanded='true']").forEach(c => {
      c.setAttribute("aria-expanded", "false");
    });
  });
});

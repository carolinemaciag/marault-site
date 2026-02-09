document.addEventListener("DOMContentLoaded", () => {
  const elements = document.querySelectorAll(
    ".service-intro, .service-row, .gain-item, .fade-in, .reveal"
  );

  if (!elements.length) return;

  if (!("IntersectionObserver" in window)) {
    elements.forEach(el => el.classList.add("is-visible"));
    return;
  }

  const observer = new IntersectionObserver(
    (entries, observer) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add("is-visible");
          observer.unobserve(entry.target);
        }
      });
    },
    {
      threshold: 0.15,
      rootMargin: "0px 0px -80px 0px"
    }
  );

  elements.forEach(el => observer.observe(el));
});











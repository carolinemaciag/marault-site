document.addEventListener("DOMContentLoaded", () => {
  const elements = document.querySelectorAll(
    ".service-intro, .service-row"
  );

  if (!elements.length) return;

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
      threshold: 0.2,
      rootMargin: "0px 0px -80px 0px"
    }
  );

  elements.forEach(el => observer.observe(el));
});


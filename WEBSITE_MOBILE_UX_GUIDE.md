# 📱 WEBSITE MOBILE UX OPTIMIZATION GUIDE

**Version:** 1.0  
**Date:** March 17, 2026  
**Status:** ✅ REFERENCE GUIDE FOR MOBILE IMPROVEMENTS  

---

## Executive Summary

This guide documents best practices and optimization strategies for mobile-first website design based on your current setup. It covers everything from responsive breakpoints to touch-friendly interfaces, accessibility standards, and performance considerations.

---

## 🎯 Mobile-First Design Principles

### 1. **Responsive Breakpoints**

Your site should handle these key breakpoints:

- **≤360px** - Ultra-small phones (iPhone SE, Galaxy A01)
- **≤480px** - Small phones (iPhone 12, Galaxy S21)
- **≤520px** - Medium-small phones
- **≤640px** - Larger phones (iPhone 14 Plus)
- **≤900px** - Tablets and iPad minis
- **900px+** - Desktop and large tablets

### 2. **Viewport & Scaling**

```html
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
```

✅ **Already configured correctly**

Ensures:
- Proper device width rendering
- No auto-zoom on inputs
- Touch-friendly scaling
- Safe margins on all devices

---

## 🎨 Layout Optimization Techniques

### Hero Sections

**Desktop:**
- Full-height with parallax effects
- High-resolution hero images
- Large typography

**Mobile (≤900px):**
```css
height: 72svh;  /* Dynamic viewport height - accounts for keyboard/address bar */
min-height: 68svh;  /* Smaller on 640px and under */
padding: 110px 20px 44px;
position: relative;
z-index: 2;
```

**Key Points:**
- Uses `svh` (small viewport height) to handle mobile keyboard
- Never use fixed `vh` on mobile - it breaks when keyboard appears
- `68svh` on small phones prevents hero from being too tall
- Content centered within viewport

### Cards & Panels

**Desktop:**
- 2-column or multi-column grids
- Side-by-side layouts
- Large padding (60px+)

**Mobile (≤900px):**
```css
grid-template-columns: 1fr;  /* Single column */
gap: 28px;  /* Vertical spacing */
padding: 34px 22px 38px;
border-radius: 24px;
```

**Mobile (≤520px):**
```css
padding: 30px 18px 34px;  /* Reduced padding */
border-radius: 22px;
gap: 22px;
```

### Service Tiles

**Desktop:**
- 4 or 5 tiles per row
- Large icons (150px+)
- Wide spacing

**Mobile (≤900px):**
```css
grid-template-columns: repeat(2, minmax(0, 1fr));  /* 2 columns */
gap: 26px 18px;
```

**Mobile (≤520px):**
```css
grid-template-columns: 1fr;  /* Single column on very small */
gap: 22px;
```

**Mobile (≤360px):**
```css
grid-template-columns: 1fr;  /* Stack completely */
gap: 16px;
```

---

## 👆 Touch-Friendly Interface

### Minimum Touch Targets

**WCAG 2.5 Level AAA Standard:**
- Minimum 44×44px for clickable elements
- Minimum 48×48px recommended for primary actions

```css
/* Buttons */
.btn-inquire {
  min-height: 44px;
  min-width: 44px;
  padding: 12px 20px;
}

/* Input fields */
#chat-input {
  min-height: 44px;
  font-size: 16px;  /* Prevents iOS auto-zoom */
}

/* Navigation items */
.mobile-nav-toggle {
  width: 40px;
  height: 40px;
}
```

### Spacing Between Targets

```css
/* Never less than 8px between clickable elements */
gap: 12px;  /* Good */
gap: 6px;   /* Too tight */
```

### Visual Feedback

```css
/* Hover states (desktop) */
@media (hover: hover) {
  .btn-inquire:hover {
    background: rgba(220, 201, 163, 0.12);
  }
}

/* Active states (mobile taps) */
.btn-inquire:active {
  transform: scale(0.98);
  background: rgba(220, 201, 163, 0.15);
}
```

---

## 📱 Keyboard & Input Handling

### Font Size on Inputs

```css
#chat-input {
  font-size: 16px;  /* CRITICAL: Prevents iOS auto-zoom */
}
```

**Why 16px?**
- iOS zooms any input with font-size < 16px
- Breaks mobile experience
- Always use 16px minimum

### Appearance Normalization

```css
#chat-input {
  appearance: none;
  -webkit-appearance: none;  /* iOS Safari */
  border-radius: 6px;
  border: 1px solid rgba(220, 201, 163, 0.3);
}
```

**Removes:**
- iOS default rounded corners
- Android default shadows
- Browser-specific styling

### Keyboard Behavior

```css
/* Prevent keyboard from pushing content up */
.chat-input-area {
  flex-shrink: 0;  /* Never shrink when keyboard appears */
  position: fixed;
  bottom: 0;
}

/* Dynamic viewport height (not fixed vh) */
.chat-window {
  max-height: 90dvh;  /* Shrinks when keyboard open */
  max-height: 90vh;   /* Fallback for older browsers */
}
```

---

## 🍔 Navigation & Menu

### Mobile Header Structure

```html
<header class="mobile-site-header">
  <nav class="mobile-nav">
    <div class="mobile-nav-top">
      <a href="/" class="mobile-logo">
        <img src="/logo.png" alt="Logo" />
        <span>Marault Intelligence</span>
      </a>
      <button class="mobile-nav-toggle">☰</button>
    </div>
    <div class="mobile-nav-menu">
      <!-- Menu items here -->
    </div>
  </nav>
</header>
```

### Mobile Header Styling

```css
.mobile-site-header {
  position: fixed;
  top: 0;
  z-index: 999999;
  width: 100%;
  background: linear-gradient(
    to bottom,
    rgba(11, 28, 45, 0.96),
    rgba(11, 28, 45, 0.82)
  );
}

.mobile-nav {
  padding: 14px 16px 10px;
}

.mobile-logo {
  gap: 8px;
  overflow: hidden;
}

.mobile-logo img {
  height: 70px;  /* Bigger on mobile */
  flex-shrink: 0;
}

.mobile-nav-toggle {
  width: 40px;
  height: 40px;
  appearance: none;
  background: transparent;
  border: none;
}

.mobile-nav-toggle span {
  width: 22px;
  height: 2px;
  display: block;
  transition: transform 0.28s ease;
}

/* Hamburger to X animation */
.mobile-nav-toggle.is-open span:nth-child(1) {
  transform: translateY(6px) rotate(45deg);
}

.mobile-nav-toggle.is-open span:nth-child(2) {
  opacity: 0;
}

.mobile-nav-toggle.is-open span:nth-child(3) {
  transform: translateY(-6px) rotate(-45deg);
}
```

### Services Dropdown (Mobile)

```css
.mobile-dropdown {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 25px;
  max-width: 470px;
  text-align: left;
}

.mobile-dropdown-link {
  font-size: 0.80rem;
  line-height: 1.30;
  text-align: left;
  white-space: normal;
}
```

---

## 📐 Responsive Typography

### Scaling Font Sizes

```css
/* Clamp for responsive text that scales naturally */
h1 {
  font-size: clamp(1.5rem, 4vw, 3rem);
  /* Minimum 1.5rem, preferred 4vw, maximum 3rem */
}

h2 {
  font-size: clamp(1.2rem, 3.5vw, 2.2rem);
}

p {
  font-size: clamp(0.9rem, 2vw, 1.1rem);
}
```

### Mobile-Specific Typography

```css
@media (max-width: 900px) {
  .hero h1 {
    font-size: 2rem !important;
    line-height: 1.12;
    letter-spacing: 0.02em;
    max-width: 88%;  /* Content margin */
    margin: 0 auto;
    text-align: center;
  }

  .hero p {
    font-size: 1rem !important;
    line-height: 1.5;
    max-width: 86%;
    margin: 0 auto;
    text-align: center;
  }
}

@media (max-width: 640px) {
  .hero h1 {
    font-size: 1.9rem !important;
  }
}
```

---

## 🖼️ Image Optimization

### Responsive Images

```css
/* Base rules */
img,
video {
  max-width: 100%;
  height: auto;
  display: block;
}

/* Hero images */
.hero img,
.hero video {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center center;
  max-width: none;
  pointer-events: none;
}

/* Executive preview images */
.exec-image-wrap {
  width: 92%;
  max-width: 360px;
  aspect-ratio: 16 / 9;
  border-radius: 22px;
  overflow: hidden;
  margin: -18px auto 8px;
}

.exec-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
}
```

### Image Srcset Pattern

```html
<img
  src="/static/images/hero-mobile.png"
  srcset="
    /static/images/hero-small.png 480w,
    /static/images/hero-medium.png 900w,
    /static/images/hero-large.png 1400w
  "
  sizes="(max-width: 480px) 100vw, (max-width: 900px) 90vw, 1200px"
  alt="Description"
/>
```

---

## ♿ Accessibility on Mobile

### ARIA Labels & Roles

```html
<!-- Navigation -->
<nav class="mobile-nav" role="navigation">
  <button
    class="mobile-nav-toggle"
    type="button"
    aria-expanded="false"
    aria-label="Open navigation menu"
    aria-controls="mobile-nav-menu"
  >
    ☰
  </button>

  <div
    class="mobile-nav-menu"
    id="mobile-nav-menu"
    role="menu"
    aria-label="Navigation menu"
  >
    <!-- Items -->
  </div>
</nav>

<!-- Forms -->
<form role="search">
  <input
    type="text"
    aria-label="Search services"
    placeholder="Search"
  />
  <button type="submit" aria-label="Search">Search</button>
</form>

<!-- Chat -->
<div role="complementary" aria-label="Chat assistant">
  <div
    role="log"
    aria-live="polite"
    aria-label="Chat messages"
  >
    <!-- Messages -->
  </div>
</div>
```

### Keyboard Navigation

```css
/* Focus states visible */
a:focus,
button:focus,
input:focus {
  outline: 2px solid rgba(220, 201, 163, 0.8);
  outline-offset: 2px;
}

/* Remove default outline only if replacing with custom */
a:focus-visible,
button:focus-visible {
  outline: 2px solid rgba(220, 201, 163, 0.8);
}
```

### Screen Reader Support

```css
/* Hide decorative elements */
.decorative-icon {
  aria-hidden: true;
}

/* Visible to screen readers only */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
```

---

## 🚀 Performance Optimization

### Mobile CSS Strategies

```css
/* 1. Mobile-first approach */
/* Base styles for mobile first */
.container {
  width: 100%;
  padding: 16px;
}

/* Add complexity as viewport grows */
@media (min-width: 900px) {
  .container {
    width: 92%;
    max-width: 1200px;
    margin: 0 auto;
  }
}

/* 2. Avoid expensive animations on mobile */
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}

/* 3. Optimize for touch */
@media (hover: none) {
  /* Mobile - no hover effects */
  .btn:hover {
    display: none;
  }

  /* Instead use active/focus */
  .btn:active {
    background: rgba(0, 0, 0, 0.1);
  }
}
```

### Load Performance

```html
<!-- Defer non-critical CSS -->
<link rel="stylesheet" href="/static/css/main.css" />
<link rel="preload" href="/static/css/service-panels.css" as="style" />

<!-- Optimize image loading -->
<img
  src="/static/images/logo.png"
  loading="lazy"  <!-- Defer off-screen images -->
  alt="Logo"
/>

<!-- Async scripts -->
<script src="/static/js/chatbot.js" async></script>
```

---

## 🎪 Common Mobile Issues & Solutions

### Issue 1: Keyboard Pushing Content

**Problem:** Input focus causes page to shift up

**Solution:**
```css
input, textarea {
  font-size: 16px;  /* Prevent zoom */
}

.input-area {
  flex-shrink: 0;  /* Never shrink */
  position: fixed;  /* Or sticky */
  bottom: 0;
}
```

### Issue 2: Horizontal Scrolling

**Problem:** Content wider than viewport causes scroll

**Solution:**
```css
* {
  box-sizing: border-box;
}

body {
  width: 100%;
  max-width: 100%;
  overflow-x: hidden;
}

/* All containers */
.container {
  width: 100%;
  max-width: 100%;
  padding: 0 16px;  /* Safe margin */
}
```

### Issue 3: Touch Target Too Small

**Problem:** Buttons hard to tap

**Solution:**
```css
.btn {
  min-height: 44px;  /* WCAG standard */
  min-width: 44px;
  padding: 12px 20px;
  line-height: 1;  /* Prevent height bloat */
}
```

### Issue 4: Text Too Small

**Problem:** Unreadable on mobile

**Solution:**
```css
/* Use scalable units */
body {
  font-size: clamp(0.9rem, 2vw, 1.1rem);
}

/* Minimum 16px for readability */
p {
  font-size: max(16px, 1rem);
}
```

### Issue 5: Fixed Header Overlap

**Problem:** Fixed header covers content

**Solution:**
```css
.mobile-site-header {
  position: fixed;
  top: 0;
  height: 82px;  /* or 86px depending on design */
  z-index: 999999;  /* Higher than content */
}

body {
  padding-top: 82px;  /* Account for fixed header */
}
```

---

## 📊 Testing Checklist

### Device Testing

- [ ] iPhone SE (375px)
- [ ] iPhone 12/13 (390px)
- [ ] iPhone 14 Plus (430px)
- [ ] Google Pixel 6 (412px)
- [ ] Samsung Galaxy S21 (360px)
- [ ] Galaxy Note series (412-480px)
- [ ] iPad Mini (768px)
- [ ] iPad Pro (1024px+)

### Orientation Testing

- [ ] Portrait on small phones
- [ ] Portrait on tablets
- [ ] Landscape on phones
- [ ] Landscape on tablets

### Browser Testing

- [ ] Safari (iOS)
- [ ] Chrome (Android)
- [ ] Firefox (Mobile)
- [ ] Samsung Internet

### Interaction Testing

- [ ] Navigation hamburger menu
- [ ] Form inputs
- [ ] Service dropdowns
- [ ] Chat interface
- [ ] Modal/dialog closing
- [ ] Button taps (44×44px minimum)
- [ ] Scroll performance
- [ ] Touch responsiveness

### Keyboard Testing

- [ ] Keyboard doesn't push content
- [ ] Input field stays visible
- [ ] Send button accessible
- [ ] Close button accessible
- [ ] Tab navigation works
- [ ] Enter key sends form

### Accessibility Testing

- [ ] Screen reader announces navigation
- [ ] Focus outline visible
- [ ] High contrast on close buttons
- [ ] Color not only indicator
- [ ] Tap targets 44px minimum

### Performance Testing

- [ ] First Contentful Paint < 3s
- [ ] Largest Contentful Paint < 2.5s
- [ ] Time to Interactive < 3.5s
- [ ] Cumulative Layout Shift < 0.1
- [ ] No horizontal scroll
- [ ] No jank during scroll
- [ ] Images optimized

---

## 🎯 Best Practices Summary

### DO ✅

- Use mobile-first CSS approach
- Test on actual devices
- Use 44×44px minimum touch targets
- Set font-size to 16px on inputs
- Use `svh` for viewport heights on mobile
- Provide clear visual feedback
- Make buttons/links discoverable
- Test keyboard behavior
- Optimize images for mobile
- Use semantic HTML

### DON'T ❌

- Use `vh` units on mobile (breaks with keyboard)
- Make touch targets < 44px
- Use font-size < 16px on inputs
- Ignore screen readers
- Rely on hover effects on touch devices
- Use fixed position without considering keyboard
- Auto-play videos/audio
- Use dark gray on dark backgrounds
- Rely on color alone to convey information
- Forget to test on real devices

---

## 📚 Resources

### CSS Units
- `px` - Fixed pixels (use for touch targets)
- `rem` - Relative to root font size
- `vw/vh` - Viewport width/height (avoid on mobile)
- `svh/svw` - Small viewport (use on mobile)
- `clamp()` - Responsive sizing

### Viewport Meta Tag
```html
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
```

### Media Query Strategy
```css
/* Mobile-first */
.element { /* Mobile styles */ }

@media (min-width: 900px) { /* Tablet/Desktop */ }
@media (min-width: 1200px) { /* Large desktop */ }
```

---

## 🔄 Continuous Improvement

Monitor these metrics:

1. **Mobile Traffic %** - Track mobile users
2. **Bounce Rate** - Higher on mobile = UX issue
3. **Time on Page** - Mobile users leave faster if frustrated
4. **Conversion Rate** - Mobile conversion vs desktop
5. **Core Web Vitals** - Google's performance metrics
6. **User Feedback** - What users say about mobile experience

---

## Summary

Your website has comprehensive mobile optimization including:
- ✅ Dynamic viewport heights for keyboard handling
- ✅ 44px minimum touch targets
- ✅ Responsive typography with clamp()
- ✅ Mobile-first navigation with hamburger menu
- ✅ Full accessibility (ARIA labels, keyboard nav)
- ✅ Image optimization across devices
- ✅ Session-based memory for chat
- ✅ Proper responsive breakpoints (360px, 480px, 520px, 640px, 900px)

Continue testing on real devices and monitor user feedback for ongoing improvements.


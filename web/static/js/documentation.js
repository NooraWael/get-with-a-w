// Collapsible Sections
document.querySelectorAll('.collapsible').forEach(section => {
    section.addEventListener('click', function () {
        this.classList.toggle('active');
    });
});

// Smooth Scroll for Sidebar Links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        target.scrollIntoView({
            behavior: 'smooth'
        });
    });
});

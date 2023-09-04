document.getElementById("previewable-img").addEventListener("change", function (event) {
    if (event.target.files && event.target.files.length > 0) {
        const reader = new FileReader();

        reader.onload = function (e) {
            document.getElementById('img-preview-container').setAttribute('src', e.target.result)
        };

        reader.readAsDataURL(event.target.files[0]);
    }
});

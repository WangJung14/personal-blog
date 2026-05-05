document.addEventListener("DOMContentLoaded", function() {
    var editorContainer = document.getElementById('editor-container');
    if (editorContainer) {
        var quill = new Quill('#editor-container', {
            theme: 'snow',
            modules: {
                toolbar: {
                    container: [
                        [{ 'header': [1, 2, 3, false] }],
                        ['bold', 'italic', 'underline'],
                        ['blockquote', 'code-block'],
                        [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                        ['link', 'image'],
                        ['clean']
                    ],
                    handlers: {
                        image: function() {
                            var input = document.createElement('input');
                            input.setAttribute('type', 'file');
                            input.setAttribute('accept', 'image/*');
                            input.click();

                            input.onchange = function() {
                                var file = input.files[0];
                                var formData = new FormData();
                                formData.append('image', file);

                                fetch('/admin/upload', {
                                    method: 'POST',
                                    body: formData
                                })
                                .then(response => response.json())
                                .then(result => {
                                    if (result.url) {
                                        var range = quill.getSelection(true);
                                        quill.insertEmbed(range.index, 'image', result.url);
                                    }
                                })
                                .catch(error => {
                                    console.error('Error:', error);
                                    alert('Image upload failed');
                                });
                            };
                        }
                    }
                }
            }
        });

        var form = document.getElementById('post-form');
        form.onsubmit = function() {
            var content = document.querySelector('input[name=content]');
            content.value = quill.root.innerHTML;
            return true;
        };
    }
});

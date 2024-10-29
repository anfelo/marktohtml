document.addEventListener("DOMContentLoaded", function() {
    const copyBtn = document.querySelector("#copy-html-button");
    const htmlContent = document.querySelector("#converted-html");

    copyBtn.addEventListener("click", () => copyToClipboard());

    const copyToClipboard = async () => {
        try {
            await navigator.clipboard.writeText(htmlContent.innerHTML);

            copyBtn.textContent = "✅";
            setTimeout(() => {
                copyBtn.textContent = "✂️";
            }, 1000);
        } catch (e) {
            console.log("Failed to copy!", e);
        }
    };
});

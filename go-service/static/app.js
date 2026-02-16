(function () {
    "use strict";

    // --- DOM elements ---
    const dropZone = document.getElementById("drop-zone");
    const fileInput = document.getElementById("file-input");
    const previewContainer = document.getElementById("preview-container");
    const previewImage = document.getElementById("preview-image");
    const submitBtn = document.getElementById("submit-btn");
    const clearBtn = document.getElementById("clear-btn");
    const jobsSection = document.getElementById("jobs-section");
    const jobList = document.getElementById("job-list");
    const resultSection = document.getElementById("result-section");
    const resultImage = document.getElementById("result-image");
    const healthDot = document.getElementById("health-indicator");
    const filterSelect = document.getElementById("filter");

    let selectedFile = null;
    const activePolls = new Map(); // jobId -> intervalId

    // --- Health check ---
    async function checkHealth() {
        try {
            const res = await fetch("/healthz");
            if (res.ok) {
                healthDot.className = "health-dot healthy";
                healthDot.title = "API is healthy";
            } else {
                throw new Error("unhealthy");
            }
        } catch {
            healthDot.className = "health-dot unhealthy";
            healthDot.title = "API is unreachable";
        }
    }

    checkHealth();
    setInterval(checkHealth, 15000);

    // --- File selection ---
    dropZone.addEventListener("click", () => fileInput.click());

    dropZone.addEventListener("dragover", (e) => {
        e.preventDefault();
        dropZone.classList.add("drag-over");
    });

    dropZone.addEventListener("dragleave", () => {
        dropZone.classList.remove("drag-over");
    });

    dropZone.addEventListener("drop", (e) => {
        e.preventDefault();
        dropZone.classList.remove("drag-over");
        const file = e.dataTransfer.files[0];
        if (file && file.type === "image/png") {
            selectFile(file);
        }
    });

    fileInput.addEventListener("change", () => {
        if (fileInput.files.length > 0) {
            selectFile(fileInput.files[0]);
        }
    });

    function selectFile(file) {
        selectedFile = file;
        const url = URL.createObjectURL(file);
        previewImage.src = url;
        previewContainer.classList.remove("hidden");
        dropZone.classList.add("hidden");
    }

    clearBtn.addEventListener("click", () => {
        selectedFile = null;
        fileInput.value = "";
        previewContainer.classList.add("hidden");
        dropZone.classList.remove("hidden");
    });

    // --- Job submission ---
    submitBtn.addEventListener("click", async () => {
        if (!selectedFile) return;

        submitBtn.disabled = true;
        submitBtn.textContent = "Submitting...";

        try {
            const formData = new FormData();
            formData.append("image", selectedFile);
            formData.append("filter", filterSelect.value);

            const res = await fetch("/jobs", { method: "POST", body: formData });
            const data = await res.json();

            addJobToList(data);

            // Reset upload UI
            selectedFile = null;
            fileInput.value = "";
            previewContainer.classList.add("hidden");
            dropZone.classList.remove("hidden");
        } catch (err) {
            console.error("Failed to submit job:", err);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = "Apply Retro Filter";
        }
    });

    // --- Job tracking ---
    function addJobToList(job) {
        jobsSection.classList.remove("hidden");

        const id = job.id || job.job_id || "unknown";
        const status = job.status || "pending";

        const li = document.createElement("li");
        li.className = "job-item";
        li.id = "job-" + id;
        li.innerHTML =
            '<span class="job-id">' + id + "</span>" +
            '<span class="job-status status-' + status + '">' + status + "</span>";
        jobList.prepend(li);

        // Start polling if not in a terminal state
        if (status !== "completed" && status !== "failed") {
            startPolling(id);
        }
    }

    function startPolling(jobId) {
        if (activePolls.has(jobId)) return;

        const intervalId = setInterval(async () => {
            try {
                const res = await fetch("/jobs/" + jobId);
                if (!res.ok) return;

                const data = await res.json();
                const status = data.status || "pending";

                updateJobStatus(jobId, status);

                if (status === "completed") {
                    stopPolling(jobId);
                    showResult(data);
                } else if (status === "failed") {
                    stopPolling(jobId);
                }
            } catch (err) {
                console.error("Poll error for job " + jobId + ":", err);
            }
        }, 2000);

        activePolls.set(jobId, intervalId);
    }

    function stopPolling(jobId) {
        const intervalId = activePolls.get(jobId);
        if (intervalId) {
            clearInterval(intervalId);
            activePolls.delete(jobId);
        }
    }

    function updateJobStatus(jobId, status) {
        const li = document.getElementById("job-" + jobId);
        if (!li) return;

        const badge = li.querySelector(".job-status");
        badge.textContent = status;
        badge.className = "job-status status-" + status;
    }

    function showResult(job) {
        if (job.result_url) {
            resultSection.classList.remove("hidden");
            resultImage.src = job.result_url;
        }
    }
})();

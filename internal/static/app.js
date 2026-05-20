// app.js

const endpointList = document.getElementById("endpoint-list");
const endpointForm = document.getElementById("endpoint-form");
const message = document.getElementById("message");

async function loadEndpoints() {
  endpointList.innerHTML = "Loading...";

  try {
    const response = await fetch("/api/endpoints");
    const endpoints = await response.json();

    endpointList.innerHTML = "";

    if (endpoints.length === 0) {
      endpointList.innerHTML = "No endpoints yet.";
      return;
    }

    for (const endpoint of endpoints) {
      const div = document.createElement("div");
      div.className = "endpoint";

      div.innerHTML = `
        <p><strong>Endpoint name:</strong></p>
          <p>${endpoint.name}</p>
        <p><strong>Webhook URL:</strong></p>
          <p>${endpoint.generated_url}</p>

        <button class="send-button">Send Test Webhook</button>
        <button class="delete-button">Delete</button>
      `;

    const deleteButton = div.querySelector(".delete-button");
    deleteButton.addEventListener("click", async () => {
      try {
        const endpointID = endpoint.id || endpoint.ID;
        const response = await fetch(`/api/endpoints/${endpointID}`, {
          method: "DELETE",
        });

        if (!response.ok) {
          throw new Error("Failed to delete endpoint");
        }

        await loadEndpoints();

      } catch (error) {
        message.textContent = error.message;
      }
    });

    const sendButton = div.querySelector(".send-button");
    sendButton.addEventListener("click", async () => {
      try{
        const endpointID = endpoint.id;
        const response = await fetch(`/webhooks/${endpointID}`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Event-Type": "dashboard.test",
          },
          body: JSON.stringify({
            message: "Hello from my dashboard",
            source: "manual-test",
            sent_at: new Date().toISOString(),
          }),
        });
        if (!response.ok) {
          throw new Error("Failed to send test webhook");
        }

        message.textContent = "Test Webhook sent.";
      } catch (error) {
        message.textContent = error.message;
      }
    });

      endpointList.appendChild(div);
    }

  } catch (error) {
    endpointList.innerHTML = "Failed to load endpoints.";
  }
}

endpointForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const name = document.getElementById("name").value;
  const targetUrl = document.getElementById("target-url").value;

  try {
    const response = await fetch("/api/endpoints", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        target_url: targetUrl,
      }),
    });

    if (!response.ok) {
      throw new Error("Failed to create endpoint");
    }

    message.textContent = "Endpoint created.";

    endpointForm.reset();

    await loadEndpoints();

  } catch (error) {
    message.textContent = error.message;
  }
});

loadEndpoints();
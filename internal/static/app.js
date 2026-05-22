// app.js

const endpointList = document.getElementById("endpoint-list");
const endpointForm = document.getElementById("endpoint-form");
const message = document.getElementById("message");
const deliveryList = document.getElementById("delivery-list")
const eventList = document.getElementById("event-list")

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
        <p><strong>Description:</strong></p>
          <p>${endpoint.description}</p>

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
        await loadEvents();
        await loadDeliveries();

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

        const statusMessage = document.createElement("p");
        statusMessage.textContent = `Test Webhook sent!`;
        statusMessage.className = "success-message";
        div.appendChild(statusMessage)
        setTimeout(() => {
          statusMessage.remove();
        }, 3000);

        await loadEvents();
        await loadDeliveries();


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
  const description = document.getElementById("description").value;

  try {
    const response = await fetch("/api/endpoints", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        target_url: targetUrl,
        description: description,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      message.textContent = errorData.error;
      throw new Error(errorData.error);
    }

    message.textContent = "Endpoint created.";

    endpointForm.reset();

    await loadEndpoints();

  } catch (error) {
    message.textContent = error.message;
  }
});

async function loadEvents() {
  eventList.innerHTML = "Loading...";

  try {
    const response = await fetch("/api/events");
    const events = await response.json();

    eventList.innerHTML = "";

    if (events.length === 0) {
      eventList.innerHTML = "No events yet.";
      return;
    }

    for (const event of events) {
      const div = document.createElement("div");
      div.className = "endpoint";

      div.innerHTML = `
        <p><strong>Endpoint name:</strong> ${event.endpoint_name}</p>
        <p><strong>Event type:</strong> ${event.event_type}</p>
        <p><strong>Received:</strong> ${event.received_at}</p>
      `;

      eventList.appendChild(div);
    }
  } catch (error) {
    eventList.innerHTML = "Failed to load events.";
  }
}

async function loadDeliveries() {
  deliveryList.innerHTML = "Loading...";

  try {
    const response = await fetch("/api/deliveries");
    const deliveries = await response.json();

    deliveryList.innerHTML = "";

    if (deliveries.length === 0) {
      deliveryList.innerHTML = "No deliveries yet.";
      return;
    }

    for (const delivery of deliveries) {
      const div = document.createElement("div");
      div.className = "endpoint";

      div.innerHTML = `
        <p><strong>Endpoint name:</strong> ${delivery.endpoint_name}</p>
        <p><strong>Status:</strong> ${delivery.status}</p>
        <p><strong>Status code:</strong> ${delivery.status_code ?? "---"}</p>
        <p><strong>Target:</strong> ${delivery.target_url}</p>
        <p><strong>Duration:</strong> ${delivery.delivery_duration_ms ?? "---"} ms</p>
        <p><strong>Created:</strong> ${delivery.created_at}</p>
      `;

      deliveryList.appendChild(div);
    }
  } catch (error) {
    deliveryList.innerHTML = "Failed to load deliveries.";
  }
}

loadEndpoints();
loadEvents();
loadDeliveries();
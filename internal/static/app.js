const message = document.getElementById("message");

function getToken() {
  return localStorage.getItem("token");
}

function setToken(token) {
  localStorage.setItem("token", token);
}

function removeToken() {
  localStorage.removeItem("token");
}

function redirectToDashboard() {
  window.location.href = "/dashboard.html";
}

function redirectToLogin() {
  window.location.href = "/";
}

async function authFetch(url, options = {}) {
  const token = getToken();

  return fetch(url, {
    ...options,
    headers: {
      ...(options.headers || {}),
      Authorization: `Bearer ${token}`,
    },
  });
}

/* ---------------------
   Login / Register page
------------------------ */

const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");

if (loginForm) {
  if (getToken()) {
    redirectToDashboard();
  }

  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    console.log("Login button clicked");
    
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Login failed");
      }

      const data = await response.json();
      console.log(data);

      setToken(data.token);
      redirectToDashboard();
    } catch (error) {
      message.textContent = error.message;
    }
  });
}

if (registerForm) {
  registerForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const name = document.getElementById("register-name").value;
    const email = document.getElementById("register-email").value;
    const password = document.getElementById("register-password").value;

    try {
      const response = await fetch("/api/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Registration failed");
      }

      message.textContent = "Account created. You can now log in.";
      registerForm.reset();
    } catch (error) {
      message.textContent = error.message;
    }
  });
}

/* --------------
   Dashboard page
----------------- */

const endpointList = document.getElementById("endpoint-list");
const endpointForm = document.getElementById("endpoint-form");
const deliveryList = document.getElementById("delivery-list");
const eventList = document.getElementById("event-list");
const logoutButton = document.getElementById("logout-button");

if (endpointList) {
  if (!getToken()) {
    redirectToLogin();
  }

  loadEndpoints();
  loadEvents();
  loadDeliveries();
}

if (logoutButton) {
  logoutButton.addEventListener("click", () => {
    removeToken();
    redirectToLogin();
  });
}

async function loadEndpoints() {
  endpointList.innerHTML = "Loading...";

  try {
    const response = await authFetch("/api/endpoints");

    if (!response.ok) {
      throw new Error("Failed to load endpoints");
    }

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
        <p>${endpoint.description || "---"}</p>

        <button class="send-button">Send Test Webhook</button>
        <button class="delete-button">Delete</button>
      `;

      div.querySelector(".delete-button").addEventListener("click", async () => {
        await deleteEndpoint(endpoint.id);
      });

      div.querySelector(".send-button").addEventListener("click", async () => {
        await sendTestWebhook(endpoint.id, endpoint.name, div);
      });

      endpointList.appendChild(div);
    }
  } catch (error) {
    endpointList.innerHTML = "Failed to load endpoints.";
    message.textContent = error.message;
  }
}

endpointForm?.addEventListener("submit", async (event) => {
  event.preventDefault();

  const name = document.getElementById("name").value;
  const targetUrl = document.getElementById("target-url").value;
  const description = document.getElementById("description").value;

  try {
    const response = await authFetch("/api/endpoints", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name,
        target_url: targetUrl,
        description,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Failed to create endpoint");
    }

    const statusMessage = document.createElement("p");
    statusMessage.textContent = "Endpoint created.";
    statusMessage.className = "success-message";

    message.innerHTML = "";
    message.appendChild(statusMessage);

    setTimeout(() => {
      statusMessage.remove();
    }, 3000);

    endpointForm.reset();

    await loadEndpoints();
  } catch (error) {
    message.textContent = error.message;
  }
});

async function deleteEndpoint(endpointID) {
  try {
    const response = await authFetch(`/api/endpoints/${endpointID}`, {
      method: "DELETE",
    });

    if (!response.ok) {
      throw new Error("Failed to delete endpoint");
    }

    message.innerHTML = "";

    const statusMessage = document.createElement("p");
    statusMessage.textContent = "Endpoint deleted.";
    statusMessage.className = "success-message";

    message.appendChild(statusMessage);

    setTimeout(() => {
      statusMessage.remove();
    }, 3000);

    await loadEndpoints();
    await loadEvents();
    await loadDeliveries();
  } catch (error) {
    message.textContent = error.message;
  }
}

async function sendTestWebhook(endpointID, endpointName, div) {
  try {
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
    statusMessage.textContent = `Test webhook sent for ${endpointName}.`;
    statusMessage.className = "success-message";

    div.appendChild(statusMessage);

    setTimeout(() => {
      statusMessage.remove();
    }, 3000);

    await loadEvents();
    await loadDeliveries();
  } catch (error) {
    message.textContent = error.message;
  }
}

async function loadEvents() {
  eventList.innerHTML = "Loading...";

  try {
    const response = await authFetch("/api/events");

    if (!response.ok) {
      throw new Error("Failed to load events");
    }

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
    message.textContent = error.message;
  }
}

async function loadDeliveries() {
  deliveryList.innerHTML = "Loading...";

  try {
    const response = await authFetch("/api/deliveries");

    if (!response.ok) {
      throw new Error("Failed to load deliveries");
    }

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
    message.textContent = error.message;
  }
}
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

function redirectToAdminDashboard() {
  window.location.href = "/admin.html";
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
  loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    
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
      setToken(data.token);
      
      if (data.is_admin) {
        redirectToAdminDashboard();
      } else {
        redirectToDashboard();
      }

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
const endpointCount = document.getElementById("endpoint-count");
const deliveryCount = document.getElementById("delivery-count");
const endpointDetail = document.getElementById("endpoint-detail");

const userCount = document.getElementById("user-count");
const eventCount = document.getElementById("event-count");
const successfulDeliveryCount = document.getElementById("successful-delivery-count");
const failedDeliveryCount = document.getElementById("failed-delivery-count");

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

let selectedEndpoint = null;

async function loadEndpoints() {
  endpointList.textContent = "Loading...";

  try {
    const response = await authFetch("/api/endpoints");

    if (!response.ok) {
      throw new Error("Failed to load endpoints");
    }

    const endpoints = await response.json();

    endpointCount.textContent = endpoints.length;
    endpointList.innerHTML = "";
    endpointDetail.innerHTML = "";

    if (endpoints.length === 0) {
      endpointList.textContent = "No endpoints yet.";
      endpointDetail.textContent = "Create an endpoint to see its details here.";
      return;
    }

    endpoints.forEach((endpoint) => {
      endpointList.appendChild(createEndpointRow(endpoint));
    });

  selectedEndpoint = null;
  endpointDetail.textContent = "Select an endpoint to view details.";

  } catch (error) {
    endpointList.textContent = "Failed to load endpoints.";
    message.textContent = error.message;
  }
}

function createEndpointRow(endpoint) {
  const row = document.createElement("button");
  row.type = "button";
  row.className = "endpoint-list-row";

  if (selectedEndpoint && endpoint.id === selectedEndpoint.id) {
    row.classList.add("selected");
  }

  const name = document.createElement("span");
  name.textContent = endpoint.name;

  const status = document.createElement("span");
  status.className = "endpoint-status";
  status.textContent = endpoint.is_active ? "● Active" : "● Inactive";

  row.append(name, status);

  row.addEventListener("click", () => {
    if (selectedEndpoint?.id === endpoint.id) {
      selectedEndpoint = null;
      updateSelectedEndpointRow(null);
      endpointDetail.textContent = "Select an endpoint to view details.";
      return;
  }
    selectedEndpoint = endpoint;
    updateSelectedEndpointRow(row);
    renderEndpointDetail(endpoint);
  });

  return row;
}

function updateSelectedEndpointRow(selectedRow) {
  document.querySelectorAll(".endpoint-list-row").forEach((row) => {
    row.classList.remove("selected");
  });
  if ( selectedRow) {
    selectedRow.classList.add("selected");
  }
}

function renderEndpointDetail(endpoint) {
  endpointDetail.innerHTML = "";

  const header = createEndpointDetailHeader(endpoint);
  const url = createEndpointUrl(endpoint);
  if (endpoint.description) {
    const desctiption = createEndpointDescription(endpoint);
    endpointDetail.append(header, url, desctiption);
  } else {
    endpointDetail.append(header, url);
  }
}

function createEndpointDetailHeader(endpoint) {
  const header = document.createElement("div");
  header.className = "endpoint-detail-header";

  const title = document.createElement("h3");
  title.textContent = endpoint.name;

  const menu = createEndpointMenu(endpoint);

  header.append(title, menu);

  return header;
}

function createEndpointUrl(endpoint) {
  const url = document.createElement("div");
  url.className = "endpoint-url";
  url.textContent = endpoint.generated_url;

  return url;
}

function createEndpointDescription(endpoint) {
  const description = document.createElement("div");
  description.className = "endpoint-description";
  description.textContent = endpoint.description;

  return description;
}

function createEndpointMenu(endpoint) {
  const wrapper = document.createElement("div");
  wrapper.className = "endpoint-menu-wrapper";

  const menuButton = document.createElement("button");
  menuButton.className = "menu-button";
  menuButton.type = "button";
  menuButton.textContent = "⋮";

  const menu = document.createElement("div");
  menu.className = "endpoint-menu hidden";

  const sendButton = createMenuButton("Send test", async () => {
    await sendTestWebhook(endpoint.id, endpoint.name, endpointDetail);
  });

  const copyButton = createMenuButton("Copy URL", async () => {
    await navigator.clipboard.writeText(endpoint.generated_url);
    message.textContent = "Webhook URL copied.";
  });

  const deleteButton = createMenuButton("Delete", async () => {
    await deleteEndpoint(endpoint.id);
  });

  deleteButton.classList.add("danger");

  menuButton.addEventListener("click", (event) => {
    event.stopPropagation();
    closeEndpointMenus();
    menu.classList.toggle("hidden");
  });

  menu.append(sendButton, copyButton, deleteButton);
  wrapper.append(menuButton, menu);

  return wrapper;
}

function createMenuButton(label, action) {
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = label;

  button.addEventListener("click", async () => {
    await action();
    closeEndpointMenus();
  });

  return button;
}

function closeEndpointMenus() {
  document.querySelectorAll(".endpoint-menu").forEach((menu) => {
    menu.classList.add("hidden");
  });
}

document.addEventListener("click", closeEndpointMenus);

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
    eventCount.textContent = events.length;
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
    deliveryCount.textContent = deliveries.length;
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

/* --------------
   Admin specific
----------------- */
const userList = document.getElementById("admin-user-list");
if (userList) {
  loadUsers();
  loadAdminStats();
}
async function loadAdminStats() {
  const response = await authFetch("/admin/stats");

  if (!response.ok) {
    throw new Error("Failed to load admin stats");
  }

  const stats = await response.json();

  userCount.textContent = stats.users;
  eventCount.textContent = stats.events_received;
  successfulDeliveryCount.textContent = stats.successful_deliveries;
  failedDeliveryCount.textContent = stats.failed_deliveries;
}

async function loadUsers() {
  try {
    const response = await authFetch("/admin/users");

    if (!response.ok) {
      throw new Error("Failed to load users");
    }

    const users = await response.json()
    userList.replaceChildren();
    userCount.textContent = users.length;

    for (const user of users) {
      const row = document.createElement("tr");

      const name = document.createElement("td");
      name.textContent = user.name;

      const email = document.createElement("td");
      email.textContent = user.email;

      const role = document.createElement("td");
      const badge = document.createElement("span");
      badge.textContent = user.is_admin ? "Admin" : "User";
      badge.className = user.is_admin ? "role-admin" : "role-user";
      role.appendChild(badge)

      const created = document.createElement("td");
      created.textContent = new Date(user.created_at).toLocaleDateString();

      const actions = document.createElement("td");

      const deleteButton = document.createElement("button");
      deleteButton.textContent = "Delete";
      deleteButton.className = "delete-button";

      deleteButton.addEventListener("click", async () => {
        await deleteUser(user.id);
      });

      actions.appendChild(deleteButton);

      row.append(
        name,
        email,
        role,
        created,
        actions,
      );
      userList.appendChild(row);
    }
  } catch (error) {
    console.error(error);
    message.textContent = error.message;
  }
}

async function deleteUser(id) {
  if (!confirm("Delete this user?")) {
    return;
  }
  const response = await authFetch(`/admin/users/${id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("failed to delete users");
  }
  await loadUsers();
}
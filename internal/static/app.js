const message = document.getElementById("message");

function getAccessToken() {
  return localStorage.getItem("access_token");
}

function setAccessToken(access_token) {
  localStorage.setItem("access_token", access_token);
}

function removeAccessToken() {
  localStorage.removeItem("access_token");
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
  let response = await fetchWithAccessToken(url, options);
  if (response.status !== 401) {
    return response;
  }

  const refreshResponse = await fetch("/api/refresh", {
    method: "POST",
    credentials: "include",
  });

  if (!refreshResponse.ok) {
    removeAccessToken();
    redirectToLogin();
    return refreshResponse;
  }

  const data = await refreshResponse.json();
  setAccessToken(data.access_token);

  response = await fetchWithAccessToken(url, options);
  return response;
}

async function fetchWithAccessToken(url, options = {}) {
  const access_token = getAccessToken();
  return fetch(url, {
    ...options,
    credentials: "include",
    headers: {
      ...(options.headers || {}),
      Authorization: `Bearer ${access_token}`,
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
      setAccessToken(data.access_token);
      
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
const deadLetterCount = document.getElementById("dead-letter-count");
const retryScheduledDeliveryCount = document.getElementById("retry-scheduled-count");
const userRecentActivityList = document.getElementById("user-recent-activity-list");

if (endpointList) {
  if (!getAccessToken()) {
    redirectToLogin();
  }

  loadEndpoints();
  loadUserRecentActivity();
}

if (logoutButton) {
  logoutButton.addEventListener("click", async () => {
    console.log("log out clicked")
    try {
      await fetch("/api/revoke", {
        method: "POST",
        credentials: "include",
      });
    } catch (error) {
      console.error(error);
    } finally {
      removeAccessToken();
      redirectToLogin();
    }
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
    await loadUserRecentActivity();
  } catch (error) {
    message.textContent = error.message;
  }
}

const eventTypes = [
  "payment.success",
  "payment.failed",
  "user.created",
  "subscription.renewed",
  "order.shipped",
];

function randomEventType() {
  return eventTypes[Math.floor(Math.random() * eventTypes.length)];
}

async function sendTestWebhook(endpointID, endpointName, div) {
  try {
    const response = await authFetch(`/webhooks/${endpointID}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Event-Type": randomEventType(),
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

    await loadUserRecentActivity();
  } catch (error) {
    message.textContent = error.message;
  }
}

function createInfoRow(label, value) {
  const row = document.createElement("p");

  const strong = document.createElement("strong");
  strong.textContent = `${label}: `;

  row.append(strong, document.createTextNode(value));

  return row;
}

async function loadUserRecentActivity() {
  try {
    const response = await authFetch("/api/recent-activity");
    if (!response.ok) {
      throw new Error("failed to load recent activity");
    }

    const activities = await response.json();
    userRecentActivityList.replaceChildren();

    for (const activity of activities) {
      const row = document.createElement("tr");
      const time = document.createElement("td");
      time.textContent = timeAgo(activity.received_at);

      const endpoint = document.createElement("td");
      endpoint.textContent = activity.endpoint_name;

      const eventType = document.createElement("td");
      eventType.textContent = activity.event_type;

      const attempts = document.createElement("td");
      attempts.textContent = `${activity.attempt_count}/5`;

      const status = document.createElement("td");
      const badge = document.createElement("span");
      badge.textContent = activity.latest_delivery_status;
      badge.className = `status-badge status-${activity.latest_delivery_status}`;
      status.appendChild(badge);

      const code = document.createElement("td");
      code.textContent = activity.latest_delivery_status_code;

      row.append(
        time,
        endpoint,
        eventType,
        attempts,
        status,
        code,
      );
      userRecentActivityList.appendChild(row);
    }
  } catch (error) {
    console.error(error);
  }
}

/* --------------
   Admin specific
----------------- */
const userList = document.getElementById("admin-user-list");
const adminEndpointList = document.getElementById("admin-endpoint-list")
const adminRecentActivityList = document.getElementById("admin-recent-activity-list")

const deliveryDetailModal = document.getElementById("delivery-detail-modal");
const deliveryDetailContent = document.getElementById("delivery-detail-content");
const deliveryDetailClose = document.getElementById("delivery-detail-close");

deliveryDetailClose?.addEventListener("click", () => {
  deliveryDetailModal.classList.add("hidden");
});

deliveryDetailModal?.addEventListener("click", (event) => {
  if (event.target === deliveryDetailModal) {
    deliveryDetailModal.classList.add("hidden");
  }
});

if (userList) {
  loadAdminUsers();
  loadAdminStats();
  loadAdminEndpoints();
  loadAdminRecentActivity();
}

function timeAgo(dateString) {
  const seconds = Math.floor((Date.now() - new Date(dateString)) / 1000);

  if (seconds < 60) {
    return `${seconds} sec ago`;
  }

  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) {
    return `${minutes} min ago`;
  }

  const hours = Math.floor(minutes / 60);
  if (hours < 24) {
    return `${hours} hr ago`;
  }

  const days = Math.floor(hours / 24);
  return `${days} day${days > 1 ? "s" : ""} ago`;
}

async function loadAdminRecentActivity() {
  try {
    const response = await authFetch("/admin/recent-activity")
    if (!response.ok) {
      throw new Error("failed to load events");
    }

    const activities = await response.json();
    adminRecentActivityList.replaceChildren();

    for (const activity of activities) {
      const row = document.createElement("tr");
      const time = document.createElement("td");
      time.textContent = timeAgo(activity.received_at);

      const endpoint = document.createElement("td");
      endpoint.textContent = activity.endpoint_name;

      const user = document.createElement("td");
      user.textContent = activity.user_name;

      const eventType = document.createElement("td");
      eventType.textContent = activity.event_type;

      const attempts = document.createElement("td");
      attempts.textContent = `${activity.attempt_count}/5`;

      const status = document.createElement("td");
      const badge = document.createElement("span");
      badge.textContent = activity.latest_delivery_status;
      badge.className = `status-badge status-${activity.latest_delivery_status}`;
      status.appendChild(badge);

      const code = document.createElement("td");
      code.textContent = activity.latest_delivery_status_code;

      const actions = document.createElement("td");
      actions.className = "activity-actions";
      actions.appendChild(createActivityMenu(activity));

      row.append(
        time,
        endpoint,
        user,
        eventType,
        attempts,
        status,
        code,
        actions,
      );
      adminRecentActivityList.appendChild(row);
    }

  } catch(error) {
    console.error(error);
  };
}

function createActivityMenu(activity) {
  const wrapper = document.createElement("div");
  wrapper.className = "endpoint-menu-wrapper";

  const menuButton = document.createElement("button");
  menuButton.className = "menu-button";
  menuButton.type = "button";
  menuButton.textContent = "⋮";

  const menu = document.createElement("div");
  menu.className = "endpoint-menu hidden";

  const viewButton = createMenuButton("View details", async () => {
    await loadAdminDeliveryDetails(activity.delivery_id);
  });

  menu.appendChild(viewButton);

  if (activity.latest_delivery_status === "dead_letter") {
    const replayButton = createMenuButton("Replay delivery", async () => {
      await replayDeadLetter(activity.delivery_id);
    });

    replayButton.classList.add("danger");
    menu.appendChild(replayButton);
  }

  menuButton.addEventListener("click", (event) => {
    event.stopPropagation();
    closeEndpointMenus();
    menu.classList.toggle("hidden");
  });

  wrapper.append(menuButton, menu);
  return wrapper;
}

async function loadAdminDeliveryDetails(deliveryID) {
  const response = await authFetch(`/admin/deliveries/${deliveryID}`);

  if (!response.ok) {
    throw new Error("Failed to load delivery details");
  }

  const delivery = await response.json();
  renderDeliveryDetails(delivery);
}

function renderDeliveryDetails(delivery) {
  deliveryDetailContent.replaceChildren();

  deliveryDetailContent.append(
    createDetailRow("Endpoint", delivery.endpoint_name),
    createDetailRow("Status", delivery.status),
    createDetailRow("Attempts", `${delivery.attempt_count}/5`),
    createDetailRow("Last Status Code", delivery.status_code ?? "---"),
    createDetailRow("Last Error", delivery.error_message ?? "---"),
    createDetailRow("Target URL", delivery.target_url),
    createDetailRow("Next Retry", formatNullableDate(delivery.next_retry_at)),
    createDetailRow("Created", formatNullableDate(delivery.created_at)),
    createDetailRow("Last Attempt", formatNullableDate(delivery.attempted_at)),
    createDetailRow("Delivered At", formatNullableDate(delivery.delivered_at)),
  );

  deliveryDetailModal.classList.remove("hidden");
}

function createDetailRow(label, value) {
  const row = document.createElement("div");
  row.className = "detail-row";

  const labelEl = document.createElement("div");
  labelEl.className = "detail-label";
  labelEl.textContent = label;

  const valueEl = document.createElement("div");
  valueEl.className = "detail-value";
  valueEl.textContent = value ?? "---";

  row.append(labelEl, valueEl);
  return row;
}

function formatNullableDate(value) {
  if (!value) {
    return "---";
  }

  return new Date(value).toLocaleString();
}

async function replayDeadLetter(deliveryID) {
  if (!confirm("Replay this dead-letter delivery?")) {
    return;
  }

  const response = await authFetch(`/admin/dead-letter/${deliveryID}/replay`, {
    method: "POST",
  });

  if (!response.ok) {
    throw new Error("Failed to replay delivery");
  }

  await loadAdminStats();
  await loadAdminRecentActivity();
}

async function loadAdminEndpoints() {
  try {
    const response = await authFetch("/admin/endpoints");
    if (!response.ok) {
      throw new Error("failed to load endpoints");
    }

    const endpoints = await response.json();
    adminEndpointList.replaceChildren();

    for (const endpoint of endpoints) {
      const row = document.createElement("tr");
      const name = document.createElement("td");
      name.textContent = endpoint.name;

      const user = document.createElement("td");
      user.textContent = endpoint.user_name;

      const status = document.createElement("td");
      status.className = "endpoint-status";
      status.textContent = endpoint.is_active ? "● Active" : "● Inactive";

      const created = document.createElement("td");
      created.textContent = new Date(endpoint.created_at).toLocaleDateString();

      row.append(
        name,
        user,
        status,
        created,
      );
      adminEndpointList.appendChild(row);
    }
  } catch (error) {
    console.error(error);
  }
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
  deadLetterCount.textContent = stats.dead_letter;
  retryScheduledDeliveryCount.textContent = stats.retry_scheduled_deliveries
}

async function loadAdminUsers() {
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

      if (!user.is_admin) {
        const deleteButton = document.createElement("button");
        deleteButton.textContent = "Delete";
        deleteButton.className = "delete-button";
        deleteButton.addEventListener("click", async () => {
        await deleteUser(user.id);
      });
      actions.appendChild(deleteButton);
    }

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
  await loadAdminUsers();
}
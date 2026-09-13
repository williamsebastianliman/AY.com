<script lang="ts">
  import { onDestroy, onMount } from "svelte";
    import api from "../lib/api";
    import LeftSideBar from "../lib/components/LeftSideBar.svelte";
    import BurgerLeftSideBar from "../lib/components/BurgerLeftSideBar.svelte";
    type Message = {
      chat_id: string;
      created_at: string;
      group_id: string;
      message: string;
      sender_id: string;
    };

    let currentMessages: Message[] = [];
    type Chat = {
      id: number;
      name: string;
      avatar: string;
      lastMessage: string;
      time: string;
      unread: number;
      isGroup?: boolean;
      messages: Message[];
    };

    let pollingInterval: any = null;

    function stopPollingMessages() {
      if (pollingInterval) {
        clearInterval(pollingInterval);
        pollingInterval = null;
      }
    }
    function startPollingMessages() {
      stopPollingMessages();
      if (!currentChat) return;
      pollingInterval = setInterval(async () => {
        await getChatsByGroup(currentChat.group_id);
      }, 500);
    }

    type User = { id: number; username: string; name: string; profile_picture_id: string, img_url: string };
    let chats: Chat[] = [
      {
        id: 1,
        name: "John Doe",
        avatar: "JD",
        lastMessage: "Hey! How are you doing?",
        time: "2m",
        unread: 2,
        messages: [
          { id: 1, text: "Hey there!", sent: false, time: "10:30 AM" },
          { id: 2, text: "Hi John! How are you?", sent: true, time: "10:32 AM" },
          { id: 3, text: "I'm doing great! Working on a new project.", sent: false, time: "10:35 AM" },
          { id: 4, text: "Hey! How are you doing?", sent: false, time: "10:45 AM" }
        ]
      },
      {
        id: 2,
        name: "Design Team",
        avatar: "DT",
        lastMessage: "Sarah: The mockups look great!",
        time: "1h",
        unread: 0,
        isGroup: true,
        messages: [
          { id: 1, text: "Let's review the new designs", sent: false, time: "9:00 AM", sender: "Mike" },
          { id: 2, text: "I'll share the latest mockups", sent: true, time: "9:05 AM" },
          { id: 3, text: "The mockups look great!", sent: false, time: "9:30 AM", sender: "Sarah" }
        ]
      },
      {
        id: 3,
        name: "Alice Smith",
        avatar: "AS",
        lastMessage: "Thanks for the help!",
        time: "3h",
        unread: 0,
        messages: [
          { id: 1, text: "Could you help me with this issue?", sent: false, time: "Yesterday" },
          { id: 2, text: "Sure! What do you need help with?", sent: true, time: "Yesterday" },
          { id: 3, text: "Thanks for the help!", sent: false, time: "Yesterday" }
        ]
      },
      {
        id: 4,
        name: "Bob Wilson",
        avatar: "BW",
        lastMessage: "See you tomorrow!",
        time: "1d",
        unread: 0,
        messages: [
          { id: 1, text: "Are we still on for the meeting tomorrow?", sent: false, time: "Yesterday" },
          { id: 2, text: "Yes, 2 PM works for me", sent: true, time: "Yesterday" },
          { id: 3, text: "Perfect! See you tomorrow!", sent: false, time: "Yesterday" }
        ]
      }
    ];

    let users: User[] = [
      // { id: 1, name: "Emma Johnson", handle: "@emmaj", avatar: "EJ" },
      // { id: 2, name: "Michael Brown", handle: "@mikeb", avatar: "MB" },
      // { id: 3, name: "Sarah Davis", handle: "@sarahd", avatar: "SD" },
      // { id: 4, name: "David Miller", handle: "@davidm", avatar: "DM" },
      // { id: 5, name: "Lisa Wilson", handle: "@lisaw", avatar: "LW" },
      // { id: 6, name: "Tom Anderson", handle: "@toma", avatar: "TA" },
      // { id: 7, name: "Jennifer Lee", handle: "@jennl", avatar: "JL" },
      // { id: 8, name: "Chris Taylor", handle: "@christ", avatar: "CT" }
    ];
    let chatInput = "";
    let currentChat: Group | null = null;
    let showModal = false;
    let isGroupMode = false;
    let selectedUsers: number[] = [];
    async function selectChat(chat: Group) {
      currentMessages = [];
      currentChat = chat;
      await getChatsByGroup(currentChat.group_id);
      startPollingMessages();
      if (isMobile) mobilePage = "chatMain";
      console.log("msg: ",currentMessages);
    }
    onDestroy(() => {
      stopPollingMessages();
    });
    function backToChatList() {
      mobilePage = "chatList";
      stopPollingMessages();
    }
    function openModal() {
      showModal = true;
      isGroupMode = false;
      selectedUsers = [];
    }
    function closeModal() {
      showModal = false;
      isGroupMode = false;
    }
    function toggleGroupMode() {
      isGroupMode = !isGroupMode;
      selectedUsers.clear();
    }
    async function startChat() {
      if (selectedUsers.length === 0) return;

      const selectedList = selectedUsers.map(id => users.find(u => u.id === id)!);

      if (isGroupMode && selectedList.length >= 2) {
        try {
          const memberIds = selectedList.map(u => u.id);
          memberIds.push(myUserId);
          await createGroup(memberIds, false);
        } catch (err) {
          alert("Failed to create group.");
          console.error(err);
        }
      } else if (selectedList.length === 1) {
        const u = selectedList[0];
        const chat = chats.find(c => c.name === u.name && !c.isGroup);
        if (!chat) {
          try {
            await createConversation(myUserId, u.id);
          } catch (err) {
            alert("Failed to start conversation.");
            console.error(err);
            return;
          }
        }
        currentChat = chat;
      }
      closeModal();
    }


    type Group = {
      group_id: string;
      is_private: boolean;
      members: string[];
      members_name: string[];
      members_profile_id: string[];
      members_public_url: string[];
      private_profile_id: string;
      private_profile_url: string;
      private_profile_name: string;
    };
    let groups: Group[] = [];

    async function handleSendMessage() {
      if (!chatInput.trim() || !currentChat) return;

      try {
        const predictRes = await api.post("/ml/predict", { text: chatInput.trim() });
        if (predictRes.data.prediction === -1) {
          alert("Bad words detected. Message not sent.");
          return;
        }
      } catch (err) {
        console.error(err as string);
        alert("Failed to check message for bad words.");
        return;
      }

      const payload = {
        group_id: currentChat.group_id,
        sender_id: myUserId,
        message: chatInput.trim()
      };

      try {
        await api.post("/chat", payload);
        chatInput = "";
      } catch (err) {
        console.error("Failed to send message:", err);
      }
    }


    function deleteConversation(id: number) {
      if (!window.confirm("Are you sure you want to delete this conversation?")) return;
      chats = chats.filter(c => c.id !== id);
      if (currentChat?.id === id) currentChat = null;
    }

    export async function fetchAllUsers(): Promise<User[]> {
      const res = await api.get("/users");
      const users: User[] = res.data;

      const enriched = await Promise.all(users.map(async (u) => {
        let img_url = "";
        if (u.profile_picture_id) {
          try {
            const mediaRes = await api.post("/media/get-media", { id: u.profile_picture_id });
            img_url = mediaRes.data.public_url || "";
          } catch {
            img_url = "";
          }
        }
        return { ...u, img_url };
      }));

      return enriched;
    }

    async function getAllGroups() {
      const rawGroups = await fetchGroupsByUser(myUserId);

      groups = await Promise.all(
        rawGroups.groups.map(async (g) => {
          g.is_private = !!g.is_private;

          let memberIds: string[] = [];
          try {
            const membersRes = await api.get(`/group/members/${g.group_id}`);
            memberIds = (membersRes.data.members || []).map((m) => m.member_id);
          } catch {
            memberIds = [];
          }

          const names: string[] = [];
          const profileIds: string[] = [];
          const publicUrls: string[] = [];

          await Promise.all(
            memberIds.map(async (memberId: string) => {
              try {
                const userRes = await api.get(`/user/${memberId}`);
                const user = userRes.data;
                names.push(user.name);
                profileIds.push(user.profile_picture_id);

                let url = "";
                if (user.profile_picture_id) {
                  try {
                    const mediaRes = await api.post("/media/get-media", { id: user.profile_picture_id });
                    url = mediaRes.data.public_url || "";
                  } catch {
                    url = "";
                  }
                }
                publicUrls.push(url);
              } catch {
                names.push("Unknown");
                profileIds.push("");
                publicUrls.push("");
              }
            })
          );

          let private_profile_id = "";
          let private_profile_url = "";
          let private_profile_name = "";
          if (g.is_private && memberIds.length === 2) {
            const otherId = memberIds.find((id) => id !== myUserId);
            if (otherId) {
              try {
                const userRes = await api.get(`/user/${otherId}`);
                const user = userRes.data;
                private_profile_id = user.profile_picture_id || "";
                private_profile_name = user.name || "";
                if (private_profile_id) {
                  try {
                    const mediaRes = await api.post("/media/get-media", { id: private_profile_id });
                    private_profile_url = mediaRes.data.public_url || "";
                  } catch {
                    private_profile_url = "";
                  }
                }
              } catch {
                private_profile_id = "";
                private_profile_url = "";
                private_profile_name = "";
              }
            }
          }

          return {
            group_id: g.group_id,
            is_private: g.is_private,
            members: memberIds,
            members_name: names,
            members_profile_id: profileIds,
            members_public_url: publicUrls,
            private_profile_id,
            private_profile_url,
            private_profile_name,
          };
        })
      );
      console.log("ag: ", groups);
    }


    async function getChatsByGroup(groupID: string) {
      const res = await fetchChatsByGroup(groupID);
      currentMessages = res.chats;
    }

    export async function fetchGroupsByUser(userId: string) {
      const res = await api.get(`/user/groups/${userId}`);
      return res.data;
    }
    export async function fetchChatsByGroup(groupId: string) {
      const res = await api.get(`/group/chats/${groupId}`);
      return res.data;
    }
    export async function createChat(groupId: string, senderId: string, message: string) {
      await api.post("/chat", {
        group_id: groupId,
        sender_id: senderId,
        message,
      });
    }
    export async function createGroup(memberIds: string[], isPrivate = false) {
      const res = await api.post("/group", {
        member_ids: memberIds,
        is_private: isPrivate,
      });
      return res.data;
    }
    export async function addGroupMember(groupId: string, userId: string) {
      await api.post("/group/add-member", {
        group_id: groupId,
        user_id: userId,
      });
    }
    export async function removeGroupMember(groupId: string, userId: string) {
      await api.post("/group/remove-member", {
        group_id: groupId,
        user_id: userId,
      });
    }

    export async function createConversation(userA: string, userB: string) {
      const res = await api.post("/conversation", {
        user_a: userA,
        user_b: userB,
      });
      return res.data;
    }
    export async function deleteChat(chatId: string) {
      await api.post("/chat/delete", {
        chat_id: chatId,
      });
    }
    export async function maskChat(chatId: string) {
      await api.post("/chat/mask", {
        chat_id: chatId,
      });
    }
    $: sortedMessages = [...currentMessages].sort(
      (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );
    let myUserId: string;
    async function loadMyProfile() {
      try {
        const meRes = await api.get("/user/get-me");
        myUserId = meRes.data.id;
        console.log(myUserId);
      } catch (err) {
        console.error("Failed to load my profile:", err);
      }
    }
    function debugSelected(){
      console.log("selected: ", selectedUsers);
    }
    function formatTime(iso: string): string {
      const date = new Date(iso);
      return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    }
    onMount(async ()=>{
      await loadMyProfile();
      users = await fetchAllUsers();
      await getAllGroups();

      console.log(users);
    });
    function handleKeydown(e: KeyboardEvent) {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        handleSendMessage();
      }
    }

    let showMembersModal = false;
    let showAddMemberModal = false;

    let currentGroupMembers: User[] = [];
    let groupToEdit: Group | null = null;
    let allUsers: User[] = [];
    let addUserLoading = false;
    let memberLoading = false;

    function openMembersModal(group: Group) {
      groupToEdit = group;
      showMembersModal = true;
      fetchGroupMembers(group.group_id);
      console.log("open");
    }
    function closeMembersModal() {
      showMembersModal = false;
      groupToEdit = null;
      currentGroupMembers = [];
    }

    async function fetchGroupMembers(groupId: string) {
      memberLoading = true;
      try {
        const res = await api.get(`/group/members/${groupId}`);
        const memberIds = res.data.members?.map((m) => m.member_id) || [];
        currentGroupMembers = await Promise.all(
          memberIds.map(async (uid: string) => {
            try {
              const userRes = await api.get(`/user/${uid}`);
              const user = userRes.data;
              let img_url = "";
              if (user.profile_picture_id) {
                try {
                  const mediaRes = await api.post("/media/get-media", { id: user.profile_picture_id });
                  img_url = mediaRes.data.public_url || "";
                } catch {}
              }
              return { ...user, img_url };
            } catch {
              return { id: uid, name: "Unknown", username: "unknown", profile_picture_id: "", img_url: "" };
            }
          })
        );
      } finally {
        memberLoading = false;
      }
    }

    async function removeMemberFromGroup(userId: string) {
      console.log("gte: ", groupToEdit);
      if (!groupToEdit) return;
      await removeGroupMember(groupToEdit.group_id, userId);
      await fetchGroupMembers(groupToEdit.group_id);
    }

    function openAddMemberModal() {
      showAddMemberModal = true;
      fetchAddableUsers();
    }
    function closeAddMemberModal() {
      showAddMemberModal = false;
    }

    async function fetchAddableUsers() {
      addUserLoading = true;
      try {
        allUsers = await fetchAllUsers();
        const memberIds = new Set(currentGroupMembers.map(m => m.id));
        allUsers = allUsers.filter(u => !memberIds.has(u.id));
      } finally {
        addUserLoading = false;
      }
    }

    async function addMemberToGroup(userId: string) {
      if (!groupToEdit) return;
      await addGroupMember(groupToEdit.group_id, userId);
      await fetchGroupMembers(groupToEdit.group_id);
      closeAddMemberModal();
    }
    async function handleUnsendMessage(message: Message) {
      if (message.sender_id !== myUserId) return;

      if (!window.confirm("Unsend this message for everyone?")) return;

      try {
        await deleteChat(message.chat_id);
        currentMessages = currentMessages.filter(m => m.chat_id !== message.chat_id);
      } catch (err) {
        alert("Failed to unsend message.");
        console.error(err);
      }
    }
    let isMobile = false;
    let mobilePage: "chatList" | "chatMain" = "chatList";
    let windowWidth = window.innerWidth;
    function handleResize() {
      windowWidth = window.innerWidth;
      isMobile = windowWidth < 700;
      if (!isMobile) mobilePage = "chatList";
    }
    onMount(() => {
      handleResize();
      window.addEventListener("resize", handleResize);
      return () => window.removeEventListener("resize", handleResize);
    });

  </script>
  <div class="app">
    <aside class="left-sidebar">
    {#if myUserId && windowWidth >= 1250}
        <LeftSideBar currentUserId={myUserId} activePage="messages" />
    {/if}
    </aside>
    {#if myUserId && windowWidth < 1250}
      <BurgerLeftSideBar currentUserId={myUserId} activePage="messages" />
    {/if}
    {#if !isMobile}
      <div class="sidebar">
        <div class="sidebar-header">
          <h1 class="sidebar-title">Messages</h1>
          <button type="button" class="new-message-btn" on:click={openModal}>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 8px;">
              <path d="M1.998 5.5c0-1.381 1.119-2.5 2.5-2.5h15c1.381 0 2.5 1.119 2.5 2.5v13c0 1.381-1.119 2.5-2.5 2.5h-15c-1.381 0-2.5-1.119-2.5-2.5v-13zm2.5-.5c-.276 0-.5.224-.5.5v2.764l8 3.638 8-3.636V5.5c0-.276-.224-.5-.5-.5h-15zm15.5 5.463l-8 3.636-8-3.638V18.5c0 .276.224.5.5.5h15c.276 0 .5-.224.5-.5v-8.037z" />
            </svg>
            New message
          </button>
        </div>
        <div class="chat-list">
          {#each groups as group (group.group_id)}
            {#if !group.is_private}
              <div
                class="chat-item {currentChat?.id === group.group_id ? "active" : ""}"
                on:click={() => selectChat(group)}
              >
                <div class="chat-avatar group-avatar">
                  GC
                </div>
                <div class="chat-content">
                  <div class="chat-header">
                    <span class="chat-name">
                      {group.members_name.filter(n => n !== "").join(", ")}
                    </span>
                  </div>
                </div>
              </div>
            {:else if group.is_private}
              <div
                class="chat-item {currentChat?.id === group.group_id ? "active" : ""}"
                on:click={() => selectChat(group)}
              >
                <div class="chat-avatar private-avatar">
                  {#if group.private_profile_url}
                    <img
                      class="avatar-img"
                      src={group.private_profile_url}
                      alt={group.private_profile_name}
                      width="100"
                      height="100"
                    />
                  {:else}
                    {#if group.private_profile_name}
                      {group.private_profile_name.split(" ").map(w => w[0]).join("").toUpperCase()}
                    {:else}
                      PC
                    {/if}
                  {/if}
                </div>
                <div class="chat-content">
                  <div class="chat-header">
                    <span class="chat-name">{group.private_profile_name}</span>
                  </div>
                </div>
              </div>
            {/if}
          {/each}
        </div>
      </div>
      <div class="chat-main">
        {#if !currentChat}
          <div class="empty-state">
            <h2>Select a message</h2>
            <p>Choose from your existing conversations, start a new one, or just keep swimming.</p>
          </div>
        {:else}
          <div class="chat-header-main">
            <div class="chat-header-info" style="cursor:pointer" on:click={() => currentChat.is_private == false && openMembersModal(currentChat)}>
              {#if currentChat.is_private}
                <div class="chat-header-avatar">
                  {#if currentChat.private_profile_url}
                    <img
                      class="avatar-img"
                      src={currentChat.private_profile_url}
                      alt={currentChat.private_profile_name}
                      width="40"
                      height="40"
                    />
                  {:else if currentChat.private_profile_name}
                    {currentChat.private_profile_name}
                  {:else}
                    PC
                  {/if}
                </div>
                <div class="chat-header-name">{currentChat.private_profile_name}</div>
              {:else}
                <div class="chat-header-avatar group-avatar">GC</div>
                <div class="chat-header-name">
                  {#if currentChat.members_name}
                    {currentChat.members_name.join(", ")}
                  {:else}
                    Group Chat
                  {/if}
                </div>
              {/if}
            </div>
            <div class="chat-options">
              <button type="button" class="option-btn" on:click={() => deleteConversation(currentChat.id)}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M16 6V4.5C16 3.11929 14.8807 2 13.5 2H10.5C9.11929 2 8 3.11929 8 4.5V6H3C2.44772 6 2 6.44772 2 7C2 7.55228 2.44772 8 3 8H4V19C4 20.1046 4.89543 21 6 21H18C19.1046 21 20 20.1046 20 19V8H21C21.5523 8 22 7.55228 22 7C22 6.44772 21.5523 6 21 6H16ZM10 4.5C10 4.22386 10.2239 4 10.5 4H13.5C13.7761 4 14 4.22386 14 4.5V6H10V4.5ZM6 8H18V19H6V8ZM9 10C9 9.44772 9.44772 9 10 9C10.5523 9 11 9.44772 11 10V17C11 17.5523 10.5523 18 10 18C9.44772 18 9 17.5523 9 17V10ZM13 10C13 9.44772 13.4477 9 14 9C14.5523 9 15 9.44772 15 10V17C15 17.5523 14.5523 18 14 18C13.4477 18 13 17.5523 13 17V10Z" />
                </svg>
              </button>
              <button type="button" class="option-btn">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 3C7.03 3 3 7.03 3 12C3 16.97 7.03 21 12 21C16.97 21 21 16.97 21 12C21 7.03 16.97 3 12 3ZM12 19C8.13 19 5 15.87 5 12C5 8.13 8.13 5 12 5C15.87 5 19 8.13 19 12C19 15.87 15.87 19 12 19ZM11 7H13V9H11V7ZM11 11H13V17H11V11Z"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="messages-container">
            {#each sortedMessages as message (message.chat_id)}
              <div class="message {message.sender_id === myUserId ? "sent" : "received"}" on:dblclick={() => handleUnsendMessage(message)}>
                <div class="message-avatar">
                  {#if message.sender_id === myUserId}
                    YOU
                  {:else}
                    {#if currentChat.is_private}
                      <img src={currentChat.private_profile_url} alt="avatar" class="avatar-img" />
                    {:else}
                      {currentChat.avatar}
                    {/if}
                  {/if}
                </div>
                <div class="message-content">
                  <div class="message-text">{message.message}</div>
                  <div class="message-time">{formatTime(message.created_at)}</div>
                </div>
              </div>
            {/each}
          </div>
          <div class="message-input-container">
            <div class="message-input-wrapper">
              <textarea
                class="message-input"
                placeholder="Start a new message"
                bind:value={chatInput}
                on:keydown={handleKeydown}
                rows={1}
              />
              <div class="input-actions">
                <button type="button" class="input-btn">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <circle cx="12" cy="12" r="10" fill="#222" />
                    <text x="12" y="16" text-anchor="middle" font-size="12" fill="#1d9bf0">😊</text>
                  </svg>
                </button>
                <button type="button" class="input-btn">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <rect x="4" y="4" width="16" height="16" rx="2" fill="#222" />
                    <circle cx="9" cy="10" r="2" fill="#1d9bf0" />
                    <path d="M4 18l4-5a2 2 0 0 1 3 0l5 7" stroke="#1d9bf0" stroke-width="1.5" fill="none"/>
                  </svg>
                </button>
                <button type="button" class="send-btn" on:click={handleSendMessage} disabled={!chatInput.trim()}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2.01 21L23 12L2.01 3L2 10L17 12L2 14L2.01 21Z" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        {/if}
      </div>
    {:else}
      {#if mobilePage === "chatList"}
        <div class="mobile-chat-list">
          <div class="sidebar-header" style="border-right:none">
            <h1 class="sidebar-title">Messages</h1>
            <button type="button" class="new-message-btn" on:click={openModal}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 8px;">
                <path d="M1.998 5.5c0-1.381 1.119-2.5 2.5-2.5h15c1.381 0 2.5 1.119 2.5 2.5v13c0 1.381-1.119 2.5-2.5 2.5h-15c-1.381 0-2.5-1.119-2.5-2.5v-13zm2.5-.5c-.276 0-.5.224-.5.5v2.764l8 3.638 8-3.636V5.5c0-.276-.224-.5-.5-.5h-15zm15.5 5.463l-8 3.636-8-3.638V18.5c0 .276.224.5.5.5h15c.276 0 .5-.224.5-.5v-8.037z" />
              </svg>
              New message
            </button>
          </div>
          <div class="chat-list" style="border-right:none">
            {#each groups as group (group.group_id)}
              {#if !group.is_private}
                <div
                  class="chat-item {currentChat?.id === group.group_id ? 'active' : ''}"
                  on:click={() => selectChat(group)}
                >
                  <div class="chat-avatar group-avatar">GC</div>
                  <div class="chat-content">
                    <div class="chat-header">
                      <span class="chat-name">{group.members_name.filter(n => n !== "").join(", ")}</span>
                    </div>
                  </div>
                </div>
              {:else if group.is_private}
                <div
                  class="chat-item {currentChat?.id === group.group_id ? 'active' : ''}"
                  on:click={() => selectChat(group)}
                >
                  <div class="chat-avatar private-avatar">
                    {#if group.private_profile_url}
                      <img class="avatar-img" src={group.private_profile_url} alt={group.private_profile_name} width="100" height="100" />
                    {:else}
                      {#if group.private_profile_name}
                        {group.private_profile_name.split(" ").map(w => w[0]).join("").toUpperCase()}
                      {:else}
                        PC
                      {/if}
                    {/if}
                  </div>
                  <div class="chat-content">
                    <div class="chat-header">
                      <span class="chat-name">{group.private_profile_name}</span>
                    </div>
                  </div>
                </div>
              {/if}
            {/each}
          </div>
        </div>
      {:else if mobilePage === "chatMain"}
        <div class="mobile-chat-main" style="width:100vw; min-width:0; max-width:100vw; height:100vh; display:flex; flex-direction:column; background:#000;">
          <div class="chat-header-main" style="padding-left:0;padding-right:0;">
            <button class="back-btn" on:click={backToChatList} aria-label="Back" style="font-size:22px; color:#1d9bf0; background:none; border:none; margin-right:10px; cursor:pointer; padding:8px 12px; border-radius:50%;">
              ←
            </button>
            <div class="chat-header-info" style="cursor:pointer" on:click={() => currentChat && openMembersModal(currentChat)}>
              {#if currentChat.is_private}
                <div class="chat-header-avatar">
                  {#if currentChat.private_profile_url}
                    <img class="avatar-img" src={currentChat.private_profile_url} alt={currentChat.private_profile_name} width="40" height="40" />
                  {:else if currentChat.private_profile_name}
                    {currentChat.private_profile_name}
                  {:else}
                    PC
                  {/if}
                </div>
                <div class="chat-header-name">{currentChat.private_profile_name}</div>
              {:else}
                <div class="chat-header-avatar group-avatar">GC</div>
                <div class="chat-header-name">
                  {#if currentChat.members_name}
                    {currentChat.members_name.join(", ")}
                  {:else}
                    Group Chat
                  {/if}
                </div>
              {/if}
            </div>
            <div class="chat-options">
              <button type="button" class="option-btn" on:click={() => deleteConversation(currentChat.id)}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M16 6V4.5C16 3.11929 14.8807 2 13.5 2H10.5C9.11929 2 8 3.11929 8 4.5V6H3C2.44772 6 2 6.44772 2 7C2 7.55228 2.44772 8 3 8H4V19C4 20.1046 4.89543 21 6 21H18C19.1046 21 20 20.1046 20 19V8H21C21.5523 8 22 7.55228 22 7C22 6.44772 21.5523 6 21 6H16ZM10 4.5C10 4.22386 10.2239 4 10.5 4H13.5C13.7761 4 14 4.22386 14 4.5V6H10V4.5ZM6 8H18V19H6V8ZM9 10C9 9.44772 9.44772 9 10 9C10.5523 9 11 9.44772 11 10V17C11 17.5523 10.5523 18 10 18C9.44772 18 9 17.5523 9 17V10ZM13 10C13 9.44772 13.4477 9 14 9C14.5523 9 15 9.44772 15 10V17C15 17.5523 14.5523 18 14 18C13.4477 18 13 17.5523 13 17V10Z" />
                </svg>
              </button>
              <button type="button" class="option-btn">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 3C7.03 3 3 7.03 3 12C3 16.97 7.03 21 12 21C16.97 21 21 16.97 21 12C21 7.03 16.97 3 12 3ZM12 19C8.13 19 5 15.87 5 12C5 8.13 8.13 5 12 5C15.87 5 19 8.13 19 12C19 15.87 15.87 19 12 19ZM11 7H13V9H11V7ZM11 11H13V17H11V11Z"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="messages-container" style="padding-bottom:10vw;">
            {#each sortedMessages as message (message.chat_id)}
              <div class="message {message.sender_id === myUserId ? 'sent' : 'received'}" on:dblclick={() => handleUnsendMessage(message)}>
                <div class="message-avatar">
                  {#if message.sender_id === myUserId}
                    YU
                  {:else}
                    {#if currentChat.is_private}
                      <img src={currentChat.private_profile_url} alt="avatar" class="avatar-img" />
                    {:else}
                      {currentChat.avatar}
                    {/if}
                  {/if}
                </div>
                <div class="message-content">
                  <div class="message-text">{message.message}</div>
                  <div class="message-time">{formatTime(message.created_at)}</div>
                </div>
              </div>
            {/each}
          </div>
          <div class="message-input-container" style="padding-bottom:env(safe-area-inset-bottom,0);">
            <div class="message-input-wrapper">
              <textarea
                class="message-input"
                placeholder="Start a new message"
                bind:value={chatInput}
                on:keydown={handleKeydown}
                rows={1}
              />
              <div class="input-actions">
                <button type="button" class="input-btn">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <circle cx="12" cy="12" r="10" fill="#222" />
                    <text x="12" y="16" text-anchor="middle" font-size="12" fill="#1d9bf0">😊</text>
                  </svg>
                </button>
                <button type="button" class="input-btn">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <rect x="4" y="4" width="16" height="16" rx="2" fill="#222" />
                    <circle cx="9" cy="10" r="2" fill="#1d9bf0" />
                    <path d="M4 18l4-5a2 2 0 0 1 3 0l5 7" stroke="#1d9bf0" stroke-width="1.5" fill="none"/>
                  </svg>
                </button>
                <button type="button" class="send-btn" on:click={handleSendMessage} disabled={!chatInput.trim()}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2.01 21L23 12L2.01 3L2 10L17 12L2 14L2.01 21Z" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      {/if}
    {/if}
    {#if showModal}
      <div class="modal-overlay" on:click|self={closeModal}>
        <div class="modal" on:click|stopPropagation>
          <div class="modal-header">
            <h2 class="modal-title">{isGroupMode ? "Create a group" : "New message"}</h2>
            <button type="button" class="close-btn" on:click={closeModal}>×</button>
          </div>
          <div class="modal-content">
            <button type="button" class="create-group-btn" style={isGroupMode ? "background: #202327" : ""} on:click={toggleGroupMode}>
              {isGroupMode ? "Cancel group creation" : "Create a group"}
            </button>
            {#each users as user (user.id)}
              <div class="user-item">
                <img class="user-avatar-img" src={user.img_url || "/assets/default-avatar.png"} alt={user.username} width="40" height="40" />
                <div class="user-info">
                  <div class="user-name">{user.username}</div>
                  {#if user.handle}
                    <div class="user-handle">{user.handle}</div>
                  {/if}
                </div>
                <input
                  type="checkbox"
                  value={user.id}
                  bind:group={selectedUsers}
                  on:change={debugSelected}
                />
              </div>
            {/each}
          </div>
          <div class="modal-footer">
            <span class="selected-count">{selectedUsers.length} selected</span>
            <button
              type="button"
              class="next-btn"
              on:click={startChat}
              disabled={selectedUsers.length === 0 || (isGroupMode && selectedUsers.length < 2)}
            >
              {isGroupMode ? (selectedUsers.length >= 2 ? "Create group" : "Next") : "Next"}
            </button>
          </div>
        </div>
      </div>
    {/if}
    {#if showMembersModal && groupToEdit}
      <div class="modal-overlay" on:click|self={closeMembersModal}>
        <div class="modal" style="max-width:400px" on:click|stopPropagation>
          <div class="modal-header">
            <h2 class="modal-title">Group Members</h2>
            <button type="button" class="close-btn" on:click={closeMembersModal}>×</button>
          </div>
          <div class="modal-content">
            {#if memberLoading}
              <div style="text-align:center; padding:2em;">Loading members...</div>
            {:else}
              {#each currentGroupMembers as user (user.id)}
                <div class="user-item" style="align-items:center;">
                  <img class="user-avatar-img" src={user.img_url || "/assets/default-avatar.png"} alt={user.username} width="40" height="40" />
                  <div class="user-info">
                    <div class="user-name">{user.name}</div>
                    <div class="user-handle">@{user.username}</div>
                  </div>
                  <button
                    type="button"
                    class="input-btn"
                    title="Remove"
                    on:click={() => removeMemberFromGroup(user.id)}
                    style="margin-left:auto; color:#f00"
                    disabled={user.id === myUserId}
                  >
                    Remove
                  </button>
                </div>
              {/each}
            {/if}
          </div>
          <div class="modal-footer">
            <button type="button" class="next-btn" on:click={openAddMemberModal}>Add Member</button>
          </div>
        </div>
      </div>
    {/if}
    {#if showAddMemberModal && groupToEdit}
        <div class="modal-overlay" on:click|self={closeAddMemberModal}>
          <div class="modal" style="max-width:400px" on:click|stopPropagation>
            <div class="modal-header">
              <h2 class="modal-title">Add to Group</h2>
              <button type="button" class="close-btn" on:click={closeAddMemberModal}>×</button>
            </div>
            <div class="modal-content">
              {#if addUserLoading}
                <div style="text-align:center; padding:2em;">Loading users...</div>
              {:else if allUsers.length === 0}
                <div style="text-align:center; padding:2em;">No users to add</div>
              {:else}
                {#each allUsers as user (user.id)}
                  <div class="user-item" style="align-items:center;">
                    <img class="user-avatar-img" src={user.img_url || "/assets/default-avatar.png"} alt={user.username} width="40" height="40" />
                    <div class="user-info">
                      <div class="user-name">{user.name}</div>
                      <div class="user-handle">@{user.username}</div>
                    </div>
                    <button
                      type="button"
                      class="input-btn"
                      style="margin-left:auto; color:#1d9bf0"
                      on:click={() => addMemberToGroup(user.id)}
                    >
                      Add
                    </button>
                  </div>
                {/each}
              {/if}
            </div>
          </div>
        </div>
      {/if}

  </div>
  <style>
    * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background: #000;
            color: #e7e9ea;
            height: 100vh;
            overflow: hidden;
        }

        .app {
            display: flex;
            height: 100vh;
        }

        .left-sidebar{
          z-index: 400000;
        }
        .sidebar {
            width: 280px;
            background: #000;
            border-right: 1px solid #2f3336;
            display: flex;
            flex-direction: column;
        }

        .sidebar-header {
            padding: 16px 20px;
            border-bottom: 1px solid #2f3336;
        }

        .sidebar-title {
            font-size: 20px;
            font-weight: 800;
            margin-bottom: 20px;
        }

        .new-message-btn {
            width: 100%;
            background: #1d9bf0;
            color: white;
            border: none;
            border-radius: 25px;
            padding: 12px 24px;
            font-size: 15px;
            font-weight: 700;
            cursor: pointer;
            transition: background-color 0.2s;
        }

        .new-message-btn:hover {
            background: #1a8cd8;
        }

        .search-box {
            margin: 16px 20px;
            position: relative;
        }

        .search-input {
            width: 100%;
            background: #202327;
            border: 1px solid #2f3336;
            border-radius: 25px;
            padding: 12px 16px 12px 45px;
            color: #e7e9ea;
            font-size: 15px;
        }

        .search-input:focus {
            outline: none;
            border-color: #1d9bf0;
        }

        .search-icon {
            position: absolute;
            left: 16px;
            top: 50%;
            transform: translateY(-50%);
            color: #71767b;
        }

        .chat-list {
            flex: 1;
            overflow-y: auto;
        }

        .chat-item {
            padding: 16px 20px;
            cursor: pointer;
            border-bottom: 1px solid #2f3336;
            transition: background-color 0.2s;
            position: relative;
            min-height: 10vh;
        }

        .chat-item:hover {
            background: #080808;
        }

        .chat-item.active {
            background: #16181c;
        }

        .chat-avatar {
            width: 48px;
            height: 48px;
            border-radius: 50%;
            background: #1d9bf0;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            color: white;
            float: left;
            margin-right: 12px;
        }

        .chat-content {
            margin-left: 60px;
        }

        .chat-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 4px;
        }

        .chat-name {
            font-weight: 700;
            font-size: 15px;
        }

        .chat-time {
            color: #71767b;
            font-size: 13px;
        }

        .chat-preview {
            color: #71767b;
            font-size: 15px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .unread-badge {
            position: absolute;
            right: 20px;
            top: 50%;
            transform: translateY(-50%);
            background: #1d9bf0;
            color: white;
            border-radius: 10px;
            padding: 2px 8px;
            font-size: 12px;
            font-weight: 700;
        }

        /* Main Chat Area */
        .chat-main {
            flex: 1;
            display: flex;
            flex-direction: column;
            background: #000;
        }

        .chat-header-main {
            padding: 16px 20px;
            border-bottom: 1px solid #2f3336;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .chat-header-info {
            display: flex;
            align-items: center;
        }

        .chat-header-avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            background: #1d9bf0;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            color: white;
            margin-right: 12px;
        }

        .chat-header-name {
            font-weight: 700;
            font-size: 15px;
        }

        .chat-options {
            display: flex;
            gap: 12px;
        }

        .option-btn {
            background: none;
            border: none;
            color: #71767b;
            cursor: pointer;
            padding: 8px;
            border-radius: 50%;
            transition: background-color 0.2s;
        }

        .option-btn:hover {
            background: #16181c;
            color: #e7e9ea;
        }

        .messages-container {
            flex: 1;
            overflow-y: auto;
            padding: 20px;
        }

        .message {
            margin-bottom: 16px;
            display: flex;
            align-items: flex-start;
        }

        .message.sent {
            flex-direction: row-reverse;
        }

        .message-avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            background: #1d9bf0;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            color: white;
            margin: 0 8px;
        }

        .message.sent .message-avatar {
            background: #00ba7c;
        }

        .message-content {
            max-width: 60%;
            background: #16181c;
            border-radius: 18px;
            padding: 12px 16px;
            position: relative;
        }

        .message.sent .message-content {
            background: #1d9bf0;
            color: white;
        }

        .message-text {
            font-size: 15px;
            line-height: 1.3;
        }

        .message-time {
            font-size: 12px;
            color: #71767b;
            margin-top: 4px;
        }

        .message.sent .message-time {
            color: rgba(255, 255, 255, 0.7);
        }

        .message-input-container {
            padding: 20px;
            border-top: 1px solid #2f3336;
        }

        .message-input-wrapper {
            display: flex;
            align-items: center;
            background: #202327;
            border-radius: 25px;
            padding: 12px 16px;
        }

        .message-input {
            flex: 1;
            background: none;
            border: none;
            color: #e7e9ea;
            font-size: 15px;
            outline: none;
            resize: none;
            max-height: 100px;
        }

        .message-input::placeholder {
            color: #71767b;
        }

        .input-actions {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-left: 12px;
        }

        .input-btn {
            background: none;
            border: none;
            color: #1d9bf0;
            cursor: pointer;
            padding: 4px;
        }

        .send-btn {
            background: #1d9bf0;
            color: white;
            border: none;
            border-radius: 50%;
            width: 32px;
            height: 32px;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .send-btn:disabled {
            background: #2f3336;
            cursor: not-allowed;
        }

        .empty-state {
            flex: 1;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            text-align: center;
            padding: 40px;
        }

        .empty-state h2 {
            font-size: 31px;
            font-weight: 800;
            margin-bottom: 8px;
        }

        .empty-state p {
            color: #71767b;
            font-size: 15px;
            line-height: 1.3;
        }

        /* Modal */
        .modal-overlay {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.4);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 1000;
        }

        .modal {
            background: #000;
            border: 1px solid #2f3336;
            border-radius: 16px;
            width: 600px;
            max-height: 80vh;
            overflow: hidden;
        }

        .modal-header {
            padding: 16px 20px;
            border-bottom: 1px solid #2f3336;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .modal-title {
            font-size: 20px;
            font-weight: 800;
        }

        .close-btn {
            background: none;
            border: none;
            color: #e7e9ea;
            cursor: pointer;
            font-size: 20px;
            padding: 8px;
            border-radius: 50%;
        }

        .close-btn:hover {
            background: #16181c;
        }

        .modal-content {
            max-height: 60vh;
            overflow-y: auto;
        }

        .user-item {
            padding: 16px 20px;
            cursor: pointer;
            border-bottom: 1px solid #2f3336;
            display: flex;
            align-items: center;
            transition: background-color 0.2s;
        }

        .user-item:hover {
            background: #080808;
        }

        .user-item.selected {
            background: #16181c;
        }

        .user-avatar {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            background: #1d9bf0;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            color: white;
            margin-right: 12px;
        }

        .user-info {
            flex: 1;
        }

        .user-name {
            font-weight: 700;
            font-size: 15px;
        }

        .user-handle {
            color: #71767b;
            font-size: 15px;
        }

        .checkbox {
            width: 20px;
            height: 20px;
            border: 2px solid #71767b;
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-left: 12px;
        }

        .checkbox.checked {
            background: #1d9bf0;
            border-color: #1d9bf0;
        }

        .modal-footer {
            padding: 16px 20px;
            border-top: 1px solid #2f3336;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .create-group-btn {
            margin-top: 2vh;
            margin-left: 2vw;
            background: #1d9bf0;
            color: white;
            border: none;
            border-radius: 25px;
            padding: 8px 16px;
            font-size: 14px;
            font-weight: 700;
            cursor: pointer;
            margin-bottom: 12px;
        }

        .selected-count {
            color: #71767b;
            font-size: 14px;
        }

        .next-btn {
            background: #1d9bf0;
            color: white;
            border: none;
            border-radius: 25px;
            padding: 12px 24px;
            font-size: 15px;
            font-weight: 700;
            cursor: pointer;
        }

        .next-btn:disabled {
            background: #2f3336;
            cursor: not-allowed;
        }

        /* Scrollbar Styling */
        ::-webkit-scrollbar {
            width: 6px;
        }

        ::-webkit-scrollbar-track {
            background: transparent;
        }

        ::-webkit-scrollbar-thumb {
            background: #2f3336;
            border-radius: 3px;
        }

        ::-webkit-scrollbar-thumb:hover {
            background: #4a4a4a;
        }
        .user-avatar-img {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          object-fit: cover;
          margin-right: 12px;
          background: #eee;
          flex-shrink: 0;
        }
        .group-avatar, .private-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1d9bf0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #fff;
  margin-right: 12px;
  overflow: hidden;
}
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
  display: block;
}

@media (max-width: 700px) {
  .sidebar {
    display: none;
  }
  .mobile-chat-list, .mobile-chat-main {
    width: 100vw;
    min-width: 0;
    max-width: 100vw;
    min-height: 0;
    max-height: 100vh;
    border: none;
    background: #000;
    flex: 1;
  }
  .chat-header-main {
    padding: 0 12px;
  }
  .messages-container {
    padding: 12px;
    min-height: 0;
    max-height: 85vh;
  }
  .message-input-container {
    padding: 12px;
    background: #000;
    border-top: 1px solid #2f3336;
  }
  .back-btn {
    font-size: 22px;
    color: #1d9bf0;
    background: none;
    border: none;
    margin-right: 10px;
    cursor: pointer;
    padding: 8px 12px;
    border-radius: 50%;
    transition: background 0.2s;
  }
  .back-btn:hover {
    background: #202327;
  }
}

</style>
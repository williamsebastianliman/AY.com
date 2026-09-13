<script>
  import { Router, Route } from "svelte-routing";
  import LandingPage from "./pages/LandingPage.svelte";
  import RegisterPage from "./pages/RegisterPage.svelte";
  import LoginPage from "./pages/LoginPage.svelte";
  import ChatPage from "./pages/ChatPage.svelte";
  import HomePage from "./pages/HomePage.svelte";
  import MePage from "./pages/MePage.svelte";
  import MyProfilePage from "./pages/MyProfilePage.svelte";
  import ProfilePage from "./pages/ProfilePage.svelte";
  import BookmarkPage from "./pages/BookmarkPage.svelte";
  import ThreadDetailPage from "./pages/ThreadDetailPage.svelte";
  import PremiumPage from "./pages/PremiumPage.svelte";
  import CommunityPage from "./pages/CommunityPage.svelte";
  import CreateCommunityPage from "./pages/CreateCommunityPage.svelte";
  import CommunityDetailPage from "./pages/CommunityDetailPage.svelte";
  import AdminPage from "./pages/AdminPage.svelte";
  import ForgetPassword from "./pages/ForgetPassword.svelte";
  import ExplorePage from "./pages/ExplorePage.svelte";
  import MessagesPage from "./pages/MessagesPage.svelte";
</script>

<Router>
  <Route path="/">
    <LandingPage />
  </Route>
  <Route path="/register">
    <RegisterPage />
  </Route>
    <Route path="/login">
    <LoginPage />
  </Route>
    <Route path="/me">
    <MePage />
  </Route>
  <Route path="/home">
    <HomePage />
  </Route>
  <Route path="/chat">
    <ChatPage />
  </Route>
  <Route path="/my-profile">
    <MyProfilePage />
  </Route>
  <Route path="/bookmark">
    <BookmarkPage />
  </Route>
  <Route path="/bookmark">
    <BookmarkPage />
  </Route>
  <Route path="/profile/:id" let:params>
    <ProfilePage {params} />
  </Route>
  <Route path="/profile/:id" let:params>
    <ProfilePage {params} />
  </Route>
  <Route path="/thread/:id" let:params>
    <ThreadDetailPage {params} />
  </Route>
  <Route path="/premium">
    <PremiumPage />
  </Route>
  <Route path="/community">
    <CommunityPage />
  </Route>
  <Route path="/create-community">
    <CreateCommunityPage />
  </Route>
  <Route path="/community/:id" let:params>
    <CommunityDetailPage {params} />
  </Route>
  <Route path="/admin">
    <AdminPage />
  </Route>
  <Route path="/admin">
    <AdminPage />
  </Route>
  <Route path="/forgot-password">
    <ForgetPassword />
  </Route>
  <Route path="/explore">
    <ExplorePage />
  </Route>
  <Route path="/messages">
    <MessagesPage />
  </Route>
</Router>
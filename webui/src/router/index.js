import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ProfileView from '../views/ProfileView.vue'
import LoginView from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'
import DiscoverUsers from '../views/DiscoverUsers.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', name: 'Login', component: LoginView },
		{
            path: '/profile/:profileId',
            name: 'Profile',
            component: ProfileView,
            meta: { requiresAuth: true }
        },
		{ path: '/stream', name:'Stream', component: StreamView , meta: { requiresAuth: true }},
		{ path: '/discover', name:'Discover', component: DiscoverUsers, meta: { requiresAuth: true } },
	]	
})

router.beforeEach((to, from, next) => {
	const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
	const isAuthenticated = localStorage.getItem('userId');
  
	if (requiresAuth && !isAuthenticated) {
	  next('/login');
	} else {
	  next();
	}
  });

export default router

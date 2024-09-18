import type {RouteRecordRaw} from "vue-router";
import {createRouter, createWebHistory} from "vue-router";

declare module "vue-router" {
    interface RouteMeta {
        parentName?: string
        requiresAuth: boolean
    }
}

const appRoutes: RouteRecordRaw[] = [
    {
        path: "/",
        name: "home",
        meta: {requiresAuth: false},
        component: () => import("@/layouts/HomeLayout.vue"),
        children:[
            {
                path: "/dashboard",
                name: "dashboard",
                meta: {requiresAuth: false},
                component: () => import("@/domain/dashboard/Dashboard.vue"),
            },
            {
                path: "/requests",
                name: "requests",
                meta: {requiresAuth: false},
                component: () => import("@/domain/requests/Requests.vue"),
            },
            {
                path: "/documents",
                name: "documents",
                meta: {requiresAuth: false},
                component: () => import("@/domain/documents/Documents.vue"),
            },
            {
                path: "/templates",
                name: "templates",
                meta: {requiresAuth: false},
                component: () => import("@/domain/templates/Templates.vue"),
            },
            {
                path: "/users",
                name: "users",
                meta: {requiresAuth: false},
                component: () => import("@/domain/users/Users.vue"),
            },
            {
                path: "/settings",
                name: "settings",
                meta: {requiresAuth: false},
                component: () => import("@/domain/settings/Settings.vue"),
            }
        ]
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: appRoutes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return {top: 0, behavior: "smooth"}
        }
    }
})

export default router
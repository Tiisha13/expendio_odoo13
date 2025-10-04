"use client";

import { PieChart, User, CheckCircle, LayoutDashboard, type LucideIcon } from "lucide-react";
import { useSession } from "next-auth/react";

import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import Link from "next/link";

export function NavBar({}: {}) {
  const { data: session } = useSession();

  if (!session) return null;

  const user = session.user;

  // Build navigation based on role
  const navlist: { name: string; url: string; icon: LucideIcon }[] = [
    { name: "Dashboard", url: "/dashboard", icon: LayoutDashboard },
  ];

  // Admins and managers can see Users page
  if (user.role === "admin" || user.role === "manager") {
    navlist.push({ name: "Users", url: "/users", icon: User });
  }

  // Everyone can see Expenses
  navlist.push({ name: "Expenses", url: "/expenses", icon: PieChart });

  // Managers and admins can see Approvals
  if (user.role === "admin" || user.role === "manager") {
    navlist.push({ name: "Approvals", url: "/approvals", icon: CheckCircle });
  }

  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarMenu>
        {navlist.map((item, index) => (
          <SidebarMenuItem key={index}>
            <SidebarMenuButton asChild>
              <Link href={item.url}>
                <item.icon />
                <span>{item.name}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}

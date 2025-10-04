"use client";

import { PieChart, User, type LucideIcon } from "lucide-react";

import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import Link from "next/link";

const navlist: { name: string; url: string; icon: LucideIcon }[] = [
  { name: "Users", url: "/users", icon: User },
  { name: "Expenses", url: "/expenses", icon: PieChart },
];

export function NavBar({}: {}) {
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

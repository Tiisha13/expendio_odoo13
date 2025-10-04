"use client";
import { GalleryVerticalEnd } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { ComponentProps, FormEvent, useState } from "react";
import CountrySelect from "./country-select";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

type SignupFormData = {
  first_name: string;
  last_name: string;
  company_name: string;
  email: string;
  password: string;
  confirm_password: string;
  country: string;
};

export function SignupForm({ className, ...props }: ComponentProps<"div">) {
  const [formData, setFormData] = useState<SignupFormData>({
    first_name: "",
    last_name: "",
    company_name: "",
    email: "",
    password: "",
    confirm_password: "",
    country: "",
  });

  const router = useRouter();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (
      !formData.first_name ||
      !formData.last_name ||
      !formData.email ||
      !formData.password ||
      !formData.confirm_password ||
      !formData.country ||
      !formData.company_name
    ) {
      toast.error("Please fill in all fields!");
      return;
    }

    // full email regex pattern
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailPattern.test(formData.email)) {
      toast.error("Please enter a valid email address!");
      return;
    }

    if (formData.password.length < 8) {
      toast.error("Password must be at least 8 characters long!");
      return;
    }
    if (formData.password !== formData.confirm_password) {
      toast.error("Passwords do not match!");
      return;
    }

    try {
      const response = await fetch(
        process.env.NEXT_PUBLIC_API_URL + "/auth/signup",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(formData),
        }
      );

      if (response.status !== 201) {
        const errorData = await response.json();
        toast.error(
          errorData.message || "Something went wrong. Please try again."
        );
        return;
      }
    } catch (error) {
      toast.error("Something went wrong. Please try again.");
      return;
    }

    toast.success("Signup successful! Please check your email to verify.", {
      action: {
        label: "Login",
        onClick: goToLogin,
      },
    });

    setTimeout(goToLogin, 5000);
  };

  const goToLogin = () => {
    router.push("/login");
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <form onSubmit={handleSubmit}>
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <Link
              href="/"
              className="flex flex-col items-center gap-2 font-medium"
            >
              <div className="flex size-8 items-center justify-center rounded-md">
                <GalleryVerticalEnd className="size-6" />
              </div>
              <span className="sr-only">Expensio.</span>
            </Link>
            <h1 className="text-xl font-bold">Welcome to Expensio.</h1>
            <FieldDescription>
              Already have an account? <Link href="/login">Sign in</Link>
            </FieldDescription>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor="first_name">First Name</FieldLabel>
              <Input
                id="first_name"
                type="text"
                required
                placeholder="John"
                onChange={(e) =>
                  setFormData({ ...formData, first_name: e.target.value })
                }
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="last_name">Last Name</FieldLabel>
              <Input
                id="last_name"
                type="text"
                required
                placeholder="Doe"
                onChange={(e) =>
                  setFormData({ ...formData, last_name: e.target.value })
                }
              />
            </Field>
          </div>
          <Field>
            <FieldLabel htmlFor="email">Email</FieldLabel>
            <Input
              id="email"
              type="email"
              placeholder="m@example.com"
              required
              onChange={(e) =>
                setFormData({ ...formData, email: e.target.value })
              }
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="country">Country</FieldLabel>
            <CountrySelect
              className="w-full"
              priorityOptions={["India"]}
              onChange={(value) => setFormData({ ...formData, country: value })}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="company">Company</FieldLabel>
            <Input
              id="company"
              type="text"
              required
              placeholder="Your Company"
              onChange={(e) =>
                setFormData({ ...formData, company_name: e.target.value })
              }
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="password">Password</FieldLabel>
            <Input
              id="password"
              type="password"
              required
              placeholder="••••••••"
              onChange={(e) =>
                setFormData({ ...formData, password: e.target.value })
              }
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="confirm-password">Confirm Password</FieldLabel>
            <Input
              id="confirm-password"
              type="password"
              required
              placeholder="••••••••"
              onChange={(e) =>
                setFormData({ ...formData, confirm_password: e.target.value })
              }
            />
          </Field>
          <Field>
            <Button type="submit">Sign up</Button>
          </Field>
        </FieldGroup>
      </form>
      <FieldDescription className="px-6 text-center">
        By clicking continue, you agree to our{" "}
        <Link href="/">Terms of Service</Link> and{" "}
        <Link href="/">Privacy Policy</Link>.
      </FieldDescription>
    </div>
  );
}

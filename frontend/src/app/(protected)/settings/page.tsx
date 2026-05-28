export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-3xl p-6">
      <h1 className="text-3xl font-bold">Settings</h1>

      <div className="mt-8 space-y-6">
        <section className="bg-card border-border rounded-xl border p-6">
          <h2 className="text-lg font-semibold">Appearance</h2>
          <p className="text-muted-foreground mt-2 text-sm">
            Customize theme preferences.
          </p>
        </section>

        <section className="bg-card border-border rounded-xl border p-6">
          <h2 className="text-lg font-semibold">Account</h2>
          <p className="text-muted-foreground mt-2 text-sm">
            Manage your account settings.
          </p>
        </section>
      </div>
    </div>
  );
}

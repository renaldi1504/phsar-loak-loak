import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

// Warna dari website Phsar (css/style.css).
const phsarPink = Color(0xFFF50750);
const phsarNavy = Color(0xFF102241);

void main() => runApp(const PhsarApp());

class PhsarApp extends StatelessWidget {
  const PhsarApp({super.key});

  @override
  Widget build(BuildContext context) {
    final base = ThemeData(
      colorScheme: ColorScheme.fromSeed(seedColor: phsarPink, primary: phsarPink),
      scaffoldBackgroundColor: Colors.white,
      useMaterial3: true,
    );
    return MaterialApp(
      title: 'Phsar',
      debugShowCheckedModeBanner: false,
      theme: base.copyWith(
        textTheme: GoogleFonts.dmSansTextTheme(base.textTheme).apply(bodyColor: phsarNavy, displayColor: phsarNavy),
      ),
      supportedLocales: const [Locale('en'), Locale('km'), Locale('zh')],
      home: const _Splash(),
    );
  }
}

class _Splash extends StatelessWidget {
  const _Splash();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Phsar', style: Theme.of(context).textTheme.displaySmall?.copyWith(fontWeight: FontWeight.w800)),
            const SizedBox(height: 6),
            const Text('Watch • Discover • Shop'),
            const SizedBox(height: 4),
            const Text('Marketplace Cambodia', style: TextStyle(color: phsarPink, fontWeight: FontWeight.w700)),
          ],
        ),
      ),
    );
  }
}
